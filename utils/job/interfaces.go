/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// Package job defines utilities for managing jobs.
package job

import (
	"context"
	"time"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/resource"
	"github.com/ARM-software/golang-utils/utils/collection/pagination"
	"github.com/ARM-software/golang-utils/utils/logs"
)

// Mocks are generated using `go generate ./...`
// Add interfaces to the following command for a mock to be generated
//go:generate go tool mockgen -destination=../mocks/mock_$GOPACKAGE.go -package=mocks github.com/ARM-software/embedded-development-services-client-utils/utils/$GOPACKAGE IAsynchronousJob,IJobManager

// IAsynchronousJob defines a typical asynchronous job.
type IAsynchronousJob interface {
	resource.IResource
	// GetDone returns whether a job has terminated.
	GetDone() bool
	// GetError returns whether a system error occurred.
	GetError() bool
	// GetFailure returns whether the job has failed.
	GetFailure() bool
	// GetSuccess returns whether the job has been successful.
	GetSuccess() bool
	// GetStatus returns the state the job is in. This is for information only and should not be relied upon as likely to change. Use flags for implementing a state machine.
	GetStatus() string
	// GetQueued returns whether the job is being queued and has not started just yet
	GetQueued() bool
	// HasMessages returns whether the job has messages available.
	HasMessages() bool
	// HasArtefacts returns whether the job has artefacts available.
	HasArtefacts() bool
}

// IJobManager defines a manager of asynchronous jobs
type IJobManager interface {
	// HasJobCompleted calls the services to determine whether the job has completed.
	HasJobCompleted(ctx context.Context, job IAsynchronousJob) (completed bool, err error)
	// HasJobStarted calls the services to determine whether the job has started.
	HasJobStarted(ctx context.Context, job IAsynchronousJob) (completed bool, err error)
	// WaitForJobCompletion waits for a job to complete. Similar to WaitForJobCompletionWithTimeout but with a timeout set to 5 minutes.
	WaitForJobCompletion(ctx context.Context, job IAsynchronousJob) (err error)
	// WaitForJobCompletionWithTimeout waits for a job to complete but with timeout protection.
	WaitForJobCompletionWithTimeout(ctx context.Context, job IAsynchronousJob, jobTimeout time.Duration) (err error)
	// GetMessagePaginator returns a paginator over job messages. The timeout corresponds to the time given to obtain the paginator
	GetMessagePaginator(ctx context.Context, logger logs.Loggers, job IAsynchronousJob, setupTimeout time.Duration) (pagination.IStreamPaginatorAndPageFetcher, error)
	// LogJobMessagesUntilNow logs all the job messages until now unless the loggingTimeout is reached beforehand. This is doing the same as WaitForJobCompletionWithTimeout apart from waiting for job completion.
	LogJobMessagesUntilNow(ctx context.Context, job IAsynchronousJob, loggingTimeout time.Duration) (err error)
}
