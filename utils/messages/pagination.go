package messages

import (
	"context"
	"time"

	"github.com/ARM-software/golang-utils/utils/collection/pagination"
)

type MessagesPaginatorFactory struct {
	backOffPeriod      time.Duration
	runOutTimeout      time.Duration
	fetchFirstPageFunc func(context.Context) (pagination.IStaticPageStream, error)
	fetchNextPageFunc  func(context.Context, pagination.IStaticPage) (pagination.IStaticPage, error)
	fetchFutureFunc    func(context.Context, pagination.IStaticPageStream) (pagination.IStaticPageStream, error)
}

func (f *MessagesPaginatorFactory) Create(ctx context.Context) (paginator pagination.IStreamPaginatorAndPageFetcher, err error) {
	return pagination.NewStaticPageStreamPaginator(ctx, f.runOutTimeout, f.backOffPeriod, f.fetchFirstPageFunc, f.fetchNextPageFunc, f.fetchFutureFunc)
}

func (f *MessagesPaginatorFactory) UpdateRunOutTimeout(runOutTimeOut time.Duration) *MessagesPaginatorFactory {
	f.runOutTimeout = runOutTimeOut
	return f
}

func (f *MessagesPaginatorFactory) UpdatBackOffPeriod(backoff time.Duration) *MessagesPaginatorFactory {
	f.backOffPeriod = backoff
	return f
}

// NewMessagePaginatorFactory returns a message paginator factory.
func NewMessagePaginatorFactory(runOutTimeOut, backoff time.Duration, fetchFirstPageFunc func(context.Context) (pagination.IStaticPageStream, error), fetchNextPageFunc func(context.Context, pagination.IStaticPage) (pagination.IStaticPage, error), fetchFutureFunc func(context.Context, pagination.IStaticPageStream) (pagination.IStaticPageStream, error)) *MessagesPaginatorFactory {
	return &MessagesPaginatorFactory{
		backOffPeriod:      backoff,
		runOutTimeout:      runOutTimeOut,
		fetchFirstPageFunc: fetchFirstPageFunc,
		fetchNextPageFunc:  fetchNextPageFunc,
		fetchFutureFunc:    fetchFutureFunc,
	}
}
