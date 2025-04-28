/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package job

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/api"
	"github.com/ARM-software/embedded-development-services-client-utils/utils/messages"
	"github.com/ARM-software/golang-utils/utils/collection/pagination"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/logs"
	"github.com/ARM-software/golang-utils/utils/parallelisation"
	"github.com/ARM-software/golang-utils/utils/reflection"
	"github.com/ARM-software/golang-utils/utils/retry"
)

type Manager struct {
	messageLoggerFactory         messages.MessageLoggerFactory
	messagesPaginatorFactory     messages.PaginatorFactory
	backOffPeriod                time.Duration
	fetchJobStatusFunc           func(ctx context.Context, jobName string) (IAsynchronousJob, *http.Response, error)
	fetchJobFirstMessagePageFunc func(ctx context.Context, jobName string) (pagination.IStaticPageStream, *http.Response, error)
}

func (m *Manager) FetchJobMessagesFirstPage(ctx context.Context, job IAsynchronousJob) (page pagination.IStaticPageStream, err error) {
	err = parallelisation.DetermineContextError(ctx)
	if err != nil {
		return
	}
	if job == nil {
		err = commonerrors.UndefinedVariable("job")
		return
	}
	jobName, err := job.FetchName()
	if err != nil {
		return
	}
	if reflection.IsEmpty(jobName) {
		err = commonerrors.UndefinedVariable("job identifier")
		return
	}
	page, err = api.GenericCallAndCheckSuccess[pagination.IStaticPageStream](ctx, fmt.Sprintf("could not fetch %v [%v]'s messages first page", job.FetchType(), jobName), func(fCtx context.Context) (pagination.IStaticPageStream, *http.Response, error) {
		return m.fetchJobFirstMessagePageFunc(fCtx, jobName)
	})
	return
}

func waitForJobState(ctx context.Context, logger logs.Loggers, job IAsynchronousJob, jobState string, checkStateFunc func(context.Context, IAsynchronousJob) (bool, error), timeout time.Duration) (err error) {
	err = parallelisation.DetermineContextError(ctx)
	if err != nil {
		return
	}

	subCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	retryCfg := retry.DefaultExponentialBackoffRetryPolicyConfiguration()
	retryCfg.RetryMax = int(float64(timeout.Milliseconds())/math.Max(float64(retryCfg.RetryWaitMin.Milliseconds()), 1)) + 1

	jobName, err := job.FetchName()
	if err != nil {
		return
	}
	if reflection.IsEmpty(jobName) {
		err = commonerrors.UndefinedVariable("job identifier")
		return
	}
	notStartedError := commonerrors.Newf(commonerrors.ErrCondition, "job [%v] has not reached the expected state [%v]", jobName, jobState)
	err = retry.RetryOnError(subCtx, logs.NewPlainLogrLoggerFromLoggers(logger), retryCfg, func() error {
		inState, subErr := checkStateFunc(subCtx, job)
		if subErr != nil {
			return subErr
		}
		if inState {
			return nil
		}
		return notStartedError
	}, fmt.Sprintf("Waiting for job [%v] to %v...", jobName, jobState), notStartedError)
	return
}

func (m *Manager) waitForJobToStart(ctx context.Context, logger logs.Loggers, job IAsynchronousJob, timeout time.Duration) error {
	return waitForJobState(ctx, logger, job, "start", m.HasJobStarted, timeout)
}

func (m *Manager) waitForJobToHaveMessagesAvailable(ctx context.Context, logger logs.Loggers, job IAsynchronousJob, timeout time.Duration) error {
	return waitForJobState(ctx, logger, job, "have messages", m.areThereMessages, timeout)
}

func (m *Manager) createMessagePaginator(ctx context.Context, job IAsynchronousJob) (paginator pagination.IStreamPaginatorAndPageFetcher, err error) {
	paginator, err = m.messagesPaginatorFactory.Create(ctx, func(subCtx context.Context) (pagination.IStaticPageStream, error) {
		return m.FetchJobMessagesFirstPage(subCtx, job)
	})
	return
}

func (m *Manager) GetMessagePaginator(ctx context.Context, logger logs.Loggers, job IAsynchronousJob, timeout time.Duration) (messagePaginator pagination.IStreamPaginatorAndPageFetcher, err error) {
	err = parallelisation.DetermineContextError(ctx)
	if err != nil {
		return
	}
	subCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	err = m.waitForJobToStart(subCtx, logger, job, timeout)
	if err != nil {
		return
	}
	err = m.waitForJobToHaveMessagesAvailable(subCtx, logger, job, timeout)
	if err != nil {
		return
	}
	messagePaginator, err = m.createMessagePaginator(ctx, job)
	return
}

func (m *Manager) WaitForJobCompletion(ctx context.Context, job IAsynchronousJob) error {
	return m.WaitForJobCompletionWithTimeout(ctx, job, 5*time.Minute)
}

func (m *Manager) WaitForJobCompletionWithTimeout(ctx context.Context, job IAsynchronousJob, timeout time.Duration) (err error) {
	err = parallelisation.DetermineContextError(ctx)
	if err != nil {
		return
	}
	messageLogger, err := m.messageLoggerFactory.Create(ctx)
	if err != nil {
		return
	}
	defer func() {
		if messageLogger != nil {
			_ = messageLogger.Close()
		}
	}()
	subCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	messagePaginator, err := m.GetMessagePaginator(subCtx, messageLogger, job, timeout)
	if err != nil {
		return
	}
	defer func() {
		if messagePaginator != nil {
			_ = messagePaginator.Close()
		}
	}()

	wait, gCtx := errgroup.WithContext(subCtx)
	wait.Go(func() error {
		return messageLogger.LogMessagesCollection(gCtx, messagePaginator)

	})
	wait.Go(func() error {
		return m.checkForMessageStreamExhaustion(gCtx, messagePaginator, job)
	})
	err = wait.Wait()
	if err != nil {
		messageLogger.LogError(err)
	}
	_, err = m.HasJobCompleted(subCtx, job)
	return
}

func (m *Manager) LogJobMessagesUntilNow(ctx context.Context, job IAsynchronousJob, timeout time.Duration) (err error) {
	err = parallelisation.DetermineContextError(ctx)
	if err != nil {
		return
	}
	messageLogger, err := m.messageLoggerFactory.Create(ctx)
	if err != nil {
		return
	}
	defer func() {
		if messageLogger != nil {
			_ = messageLogger.Close()
		}
	}()
	subCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	messagePaginator, err := m.GetMessagePaginator(subCtx, messageLogger, job, timeout)
	if err != nil {
		return
	}
	defer func() {
		if messagePaginator != nil {
			_ = messagePaginator.Close()
		}
	}()

	err = messagePaginator.DryUp()
	if err != nil {
		messageLogger.LogError(err)
		return
	}
	err = messageLogger.LogMessagesCollection(subCtx, messagePaginator)
	if err != nil {
		messageLogger.LogError(err)
	}
	return
}

func (m *Manager) checkForMessageStreamExhaustion(ctx context.Context, paginator pagination.IGenericStreamPaginator, job IAsynchronousJob) error {

	for {
		err := parallelisation.DetermineContextError(ctx)
		if err != nil {
			return err
		}
		completed, err := m.HasJobCompleted(ctx, job)
		if commonerrors.Any(err, commonerrors.ErrUndefined) {
			return err
		}
		if completed {
			err = paginator.DryUp()
			return err
		}
		parallelisation.SleepWithContext(ctx, m.backOffPeriod)
	}
}

func (m *Manager) areThereMessages(ctx context.Context, job IAsynchronousJob) (hasMessages bool, err error) {
	err = parallelisation.DetermineContextError(ctx)
	if err != nil {
		return
	}
	if job == nil {
		err = commonerrors.UndefinedVariable("job")
		return
	}
	if job.HasMessages() {
		hasMessages = true
		return
	}

	jobName, err := job.FetchName()
	if err != nil {
		return
	}
	if reflection.IsEmpty(jobName) {
		err = commonerrors.UndefinedVariable("job identifier")
		return
	}
	jobType := job.FetchType()
	jobStatus, err := api.GenericCallAndCheckSuccess[IAsynchronousJob](ctx, fmt.Sprintf("could not fetch %v [%v]'s status", jobType, jobName), func(fCtx context.Context) (IAsynchronousJob, *http.Response, error) {
		return m.fetchJobStatusFunc(fCtx, jobName)
	})
	if err != nil {
		return
	}
	hasMessages = jobStatus.HasMessages()
	return
}

func (m *Manager) HasJobStarted(ctx context.Context, job IAsynchronousJob) (started bool, err error) {
	err = parallelisation.DetermineContextError(ctx)
	if err != nil {
		return
	}
	if job == nil {
		err = commonerrors.UndefinedVariable("job")
		return
	}
	if job.GetDone() {
		started = true
		return
	}
	if !job.GetQueued() {
		started = true
		return
	}

	jobName, err := job.FetchName()
	if err != nil {
		return
	}
	if reflection.IsEmpty(jobName) {
		err = commonerrors.UndefinedVariable("job identifier")
		return
	}
	jobStatus, err := api.GenericCallAndCheckSuccess[IAsynchronousJob](ctx, fmt.Sprintf("could not fetch %v [%v]'s status", job.FetchType(), jobName), func(fCtx context.Context) (IAsynchronousJob, *http.Response, error) {
		return m.fetchJobStatusFunc(fCtx, jobName)
	})
	if err != nil {
		return
	}
	if jobStatus.GetDone() {
		started = true
	} else {
		started = !jobStatus.GetQueued()
	}
	return
}

func (m *Manager) HasJobCompleted(ctx context.Context, job IAsynchronousJob) (completed bool, err error) {
	err = parallelisation.DetermineContextError(ctx)
	if err != nil {
		return
	}
	if job == nil {
		err = commonerrors.UndefinedVariable("job")
		return
	}
	jobName, err := job.FetchName()
	if err != nil {
		return
	}
	if reflection.IsEmpty(jobName) {
		err = commonerrors.UndefinedVariable("job identifier")
		return
	}
	jobType := job.FetchType()
	jobStatus, err := api.GenericCallAndCheckSuccess[IAsynchronousJob](ctx, fmt.Sprintf("could not fetch %v [%v]'s status", jobType, jobName), func(fCtx context.Context) (IAsynchronousJob, *http.Response, error) {
		return m.fetchJobStatusFunc(fCtx, jobName)
	})
	if err != nil {
		return
	}
	if jobStatus.GetDone() {
		completed = true
	}
	if jobStatus.GetError() {
		err = commonerrors.Newf(commonerrors.ErrUnexpected, "%v [%v] errored: %v", jobType, jobName, jobStatus.GetStatus())
		return
	}
	if jobStatus.GetFailure() {
		err = commonerrors.Newf(commonerrors.ErrInvalid, "%v [%v] failed: %v", jobType, jobName, jobStatus.GetStatus())
		return
	}
	if jobStatus.GetSuccess() {
		return
	}
	if completed {
		err = commonerrors.Newf(commonerrors.ErrUnexpected, "%v [%v] completed but without success: %v", jobType, jobName, jobStatus.GetStatus())
		return
	}
	return
}

// NewJobManager creates a new job manager.
func NewJobManager(logger *messages.MessageLoggerFactory, backOffPeriod time.Duration,
	fetchJobStatusFunc func(ctx context.Context, jobName string) (IAsynchronousJob, *http.Response, error),
	fetchJobFirstMessagePageFunc func(ctx context.Context, jobName string) (pagination.IStaticPageStream, *http.Response, error),
	fetchNextJobMessagesPageFunc func(context.Context, pagination.IStaticPage) (pagination.IStaticPage, error),
	fetchFutureJobMessagesPageFunc func(context.Context, pagination.IStaticPageStream) (pagination.IStaticPageStream, error)) (IJobManager, error) {
	return newJobManagerFromMessageFactory(logger, backOffPeriod, fetchJobStatusFunc, fetchJobFirstMessagePageFunc, messages.NewPaginatorFactory(messages.DefaultStreamExhaustionGracePeriod, backOffPeriod, fetchNextJobMessagesPageFunc, fetchFutureJobMessagesPageFunc))
}

func newJobManagerFromMessageFactory(logger *messages.MessageLoggerFactory, backOffPeriod time.Duration,
	fetchJobStatusFunc func(ctx context.Context, jobName string) (IAsynchronousJob, *http.Response, error),
	fetchJobFirstMessagePageFunc func(ctx context.Context, jobName string) (pagination.IStaticPageStream, *http.Response, error),
	messagePaginator *messages.PaginatorFactory) (*Manager, error) {
	if logger == nil {
		return nil, commonerrors.ErrNoLogger
	}
	if messagePaginator == nil {
		return nil, commonerrors.UndefinedVariable("paginator factory")
	}
	return &Manager{
		messageLoggerFactory:         *logger,
		messagesPaginatorFactory:     *messagePaginator,
		backOffPeriod:                backOffPeriod,
		fetchJobStatusFunc:           fetchJobStatusFunc,
		fetchJobFirstMessagePageFunc: fetchJobFirstMessagePageFunc,
	}, nil
}
