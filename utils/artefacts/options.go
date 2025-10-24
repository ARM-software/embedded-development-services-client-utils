package artefacts

import (
	"github.com/ARM-software/golang-utils/utils/logs"
)

type DownloadOptions struct {
	StopOnFirstError      bool
	MaintainTreeStructure bool
	Logger                logs.Loggers
}

type DownloadOption func(*DownloadOptions)

func newDefaultDownloadOptions() *DownloadOptions {
	return &DownloadOptions{
		StopOnFirstError:      true,
		MaintainTreeStructure: false,
		Logger:                nil,
	}
}
func NewDownloadOptions(opts ...DownloadOption) (options *DownloadOptions) {
	options = newDefaultDownloadOptions()
	for _, opt := range opts {
		opt(options)
	}
	return
}

func WithStopOnFirstError(stop bool) DownloadOption {
	return func(o *DownloadOptions) {
		o.StopOnFirstError = stop
	}
}

func WithMaintainStructure(maintain bool) DownloadOption {
	return func(o *DownloadOptions) {
		o.MaintainTreeStructure = maintain
	}
}

func WithLogger(l logs.Loggers) DownloadOption {
	return func(o *DownloadOptions) {
		o.Logger = l
	}
}
