/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
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

// WithStopOnFirstError specifies whether the Arteftact manager will stop downloading artefacts if it encounters an error from one of the artefacts.
func WithStopOnFirstError(stop bool) DownloadOption {
	return func(o *DownloadOptions) {
		o.StopOnFirstError = stop
	}
}

// WithMaintainStructure specifies whether to keep the tree structure of the artefacts or not in the output directory.
func WithMaintainStructure(maintain bool) DownloadOption {
	return func(o *DownloadOptions) {
		o.MaintainTreeStructure = maintain
	}
}

// WithLogger specifies an optional logger that is used to log downloading artefacts or errors encountered while downloading.
func WithLogger(l logs.Loggers) DownloadOption {
	return func(o *DownloadOptions) {
		o.Logger = l
	}
}
