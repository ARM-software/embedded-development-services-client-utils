/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package messages

import (
	"context"
	"time"

	"github.com/ARM-software/golang-utils/utils/collection/pagination"
)

type PaginatorFactory struct {
	backOffPeriod     time.Duration
	runOutTimeout     time.Duration
	fetchNextPageFunc func(context.Context, pagination.IStaticPage) (pagination.IStaticPage, error)
	fetchFutureFunc   func(context.Context, pagination.IStaticPageStream) (pagination.IStaticPageStream, error)
}

func (f *PaginatorFactory) Create(ctx context.Context, fetchFirstPageFunc func(context.Context) (pagination.IStaticPageStream, error)) (paginator pagination.IStreamPaginatorAndPageFetcher, err error) {
	return pagination.NewStaticPageStreamPaginator(ctx, f.runOutTimeout, f.backOffPeriod, fetchFirstPageFunc, f.fetchNextPageFunc, f.fetchFutureFunc)
}

func (f *PaginatorFactory) UpdateRunOutTimeout(runOutTimeOut time.Duration) *PaginatorFactory {
	f.runOutTimeout = runOutTimeOut
	return f
}

func (f *PaginatorFactory) UpdateBackOffPeriod(backoff time.Duration) *PaginatorFactory {
	f.backOffPeriod = backoff
	return f
}

// NewPaginatorFactory returns a message paginator factory.
func NewPaginatorFactory(runOutTimeOut, backoff time.Duration, fetchNextPageFunc func(context.Context, pagination.IStaticPage) (pagination.IStaticPage, error), fetchFutureFunc func(context.Context, pagination.IStaticPageStream) (pagination.IStaticPageStream, error)) *PaginatorFactory {
	return &PaginatorFactory{
		backOffPeriod:     backoff,
		runOutTimeout:     runOutTimeOut,
		fetchNextPageFunc: fetchNextPageFunc,
		fetchFutureFunc:   fetchFutureFunc,
	}
}
