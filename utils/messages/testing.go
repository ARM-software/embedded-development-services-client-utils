/*
 * Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package messages

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/bxcodec/faker/v3"
	"go.uber.org/atomic"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/links/linkstest"
	pagination2 "github.com/ARM-software/embedded-development-services-client-utils/utils/pagination"
	"github.com/ARM-software/embedded-development-services-client/client"
	"github.com/ARM-software/golang-utils/utils/collection/pagination"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/parallelisation"
)

func newMockHalLink() (link *client.HalLinkData, err error) {
	faker.ResetUnique()
	linkP, err := linkstest.NewFakeLink()
	if err != nil {
		return
	}
	lP := *linkP
	l, ok := lP.(client.HalLinkData)
	if !ok {
		err = fmt.Errorf("%w: invalid link type [%T]", commonerrors.ErrInvalid, lP)
		return
	}
	link = &l
	return
}
func newMockHalFeedLinks(ctx context.Context, hasNext, hasFuture bool) (links *client.HalFeedLinks, err error) {
	err = parallelisation.DetermineContextError(ctx)
	if err != nil {
		return
	}
	selfP, err := linkstest.NewFakeLink()
	if err != nil {
		return
	}
	self, ok := (*selfP).(client.HalLinkData)
	if !ok {
		err = fmt.Errorf("%w: invalid link type [%T]", commonerrors.ErrInvalid, *selfP)
		return
	}
	links = client.NewHalFeedLinks(self)
	if hasNext {
		next, err := newMockHalLink()
		if err != nil {
			return nil, err
		}
		links.SetNext(*next)
	}

	if hasFuture {
		future, err := newMockHalLink()
		if err != nil {
			return nil, err
		}
		links.SetFuture(*future)
	}
	return
}

// NewMockNotificationFeedPage generates a mock message page for testing
func NewMockNotificationFeedPage(ctx context.Context, hasNext, hasFuture bool) (f *client.NotificationFeed, err error) {
	faker.ResetUnique()
	links, err := newMockHalFeedLinks(ctx, hasNext, hasFuture)
	if err != nil {
		return
	}
	messageCount := rand.Intn(50) //nolint:gosec //causes G404: Use of weak random number generator
	var messages []client.NotificationMessageObject
	for i := 0; i < messageCount; i++ {
		messages = append(messages, *client.NewNotificationMessageObject(faker.Sentence()))
	}
	f = client.NewNotificationFeed(*client.NewNullableHalFeedLinks(links), *client.NewNullablePagingMetadata(client.NewPagingMetadata(int32(messageCount), time.Now(), 50, time.Now(), 45, 100)), messages, faker.Name())
	return
}

// NewMockNotificationFeedPaginator generates a mock message paginator for testing
func NewMockNotificationFeedPaginator(ctx context.Context) (pagination.IPaginatorAndPageFetcher, error) {
	faker.ResetUnique()
	pageNumber := rand.Intn(50) //nolint:gosec //causes G404: Use of weak random number generator
	pageCount := atomic.NewInt64(0)
	return pagination.NewStaticPagePaginator(ctx, func(fctx context.Context) (pagination.IStaticPage, error) {
		firstPage, err := NewMockNotificationFeedPage(fctx, pageNumber > 0, false)
		if err != nil {
			return nil, err
		}
		return pagination2.ToPage(firstPage), nil
	}, func(gCtx context.Context, page pagination.IStaticPage) (pagination.IStaticPage, error) {
		pageCount.Inc()
		newPage, err := NewMockNotificationFeedPage(gCtx, int64(pageNumber) != pageCount.Load(), false)
		if err != nil {
			return nil, err
		}
		return pagination2.ToPage(newPage), nil
	})
}

// NewMockNotificationFeedStreamPaginator generates a mock message stream paginator for testing
func NewMockNotificationFeedStreamPaginator(ctx context.Context) (pagination.IStreamPaginatorAndPageFetcher, error) {
	pageNumber := rand.Intn(50) //nolint:gosec //causes G404: Use of weak random number generator
	pageCount := atomic.NewInt64(0)
	return pagination.NewStaticPageStreamPaginator(ctx, SleepAtEnd, SleepBetweenMessages, func(fctx context.Context) (pagination.IStaticPageStream, error) {
		firstPage, err := NewMockNotificationFeedPage(fctx, pageNumber > 0, false)
		if err != nil {
			return nil, err
		}
		return pagination2.ToStream(firstPage), nil
	}, func(gCtx context.Context, page pagination.IStaticPage) (pagination.IStaticPage, error) {
		pageCount.Inc()
		newPage, err := NewMockNotificationFeedPage(gCtx, int64(pageNumber) != pageCount.Load(), false)
		if err != nil {
			return nil, err
		}
		return pagination2.ToPage(newPage), nil

	}, func(ctx context.Context, stream pagination.IStaticPageStream) (pagination.IStaticPageStream, error) {
		return nil, nil
	})
}
