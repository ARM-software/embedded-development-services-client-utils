/*
 * Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package job

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/api"
	"github.com/ARM-software/embedded-development-services-client-utils/utils/messages"
	"github.com/ARM-software/golang-utils/utils/collection/pagination"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/parallelisation"
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
		err = fmt.Errorf("%w: missing job", commonerrors.ErrUndefined)
		return
	}
	jobName, err := job.FetchName()
	if err != nil {
		return
	}
	jobType := job.FetchType()
	page, resp, apierr := m.fetchJobFirstMessagePageFunc(ctx, jobName)
	if resp != nil {
		_ = resp.Body.Close()
	}
	err = api.CheckAPICallSuccess(ctx, fmt.Sprintf("could not fetch %v [%v]'s messages first page", jobType, jobName), resp, apierr)
	if err != nil {
		return
	}
	return
}

func (m *Manager) WaitForJobCompletion(ctx context.Context, job IAsynchronousJob) (err error) {
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
	messagePaginator, err := m.messagesPaginatorFactory.Create(ctx, func(subCtx context.Context) (pagination.IStaticPageStream, error) {
		return m.FetchJobMessagesFirstPage(subCtx, job)
	})
	if err != nil {
		return
	}
	defer func() {
		if messagePaginator != nil {
			_ = messagePaginator.Close()
		}
	}()

	wait, gCtx := errgroup.WithContext(ctx)
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
	_, err = m.HasJobCompleted(ctx, job)
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

func (m *Manager) HasJobCompleted(ctx context.Context, job IAsynchronousJob) (completed bool, err error) {
	err = parallelisation.DetermineContextError(ctx)
	if err != nil {
		return
	}
	if job == nil {
		err = fmt.Errorf("%w: missing job", commonerrors.ErrUndefined)
		return
	}
	jobName, err := job.FetchName()
	if err != nil {
		return
	}
	jobType := job.FetchType()
	jobStatus, resp, apierr := m.fetchJobStatusFunc(ctx, jobName)
	if resp != nil {
		_ = resp.Body.Close()
	}
	err = api.CheckAPICallSuccess(ctx, fmt.Sprintf("could not fetch %v [%v]'s status", jobType, jobName), resp, apierr)
	if err != nil {
		return
	}
	if jobStatus.GetDone() {
		completed = true
	}
	if jobStatus.GetError() {
		err = fmt.Errorf("%w: %v [%v] errored: %v", commonerrors.ErrUnexpected, jobType, jobName, jobStatus.GetStatus())
		return
	}
	if jobStatus.GetFailure() {
		err = fmt.Errorf("%w: %v [%v] failed: %v", commonerrors.ErrInvalid, jobType, jobName, jobStatus.GetStatus())
		return
	}
	if jobStatus.GetSuccess() {
		return
	}
	if completed {
		err = fmt.Errorf("%w: %v [%v] completed but without success: %v", commonerrors.ErrUnexpected, jobType, jobName, jobStatus.GetStatus())
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
		return nil, fmt.Errorf("%w: missing paginator factory", commonerrors.ErrUndefined)
	}
	return &Manager{
		messageLoggerFactory:         *logger,
		messagesPaginatorFactory:     *messagePaginator,
		backOffPeriod:                backOffPeriod,
		fetchJobStatusFunc:           fetchJobStatusFunc,
		fetchJobFirstMessagePageFunc: fetchJobFirstMessagePageFunc,
	}, nil
}
