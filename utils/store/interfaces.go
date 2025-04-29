/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// Package store provides utilities for storing data or state in a persistent fashion.
package store

import (
	"context"
	"io"

	"github.com/ARM-software/golang-utils/utils/filesystem"
)

// Mocks are generated using `go generate ./...`
// Add interfaces to the following command for a mock to be generated
//go:generate go tool mockgen -destination=../mocks/mock_$GOPACKAGE.go -package=mocks github.com/ARM-software/embedded-development-services-client-utils/utils/$GOPACKAGE IStore

type IStore interface {
	io.Closer
	// GetPath returns the store path
	GetPath() string
	// SetPath sets the store path
	SetPath(path string) error
	// Create sets up the store
	Create(ctx context.Context) error
	// Exists returns whether the store exists or not
	Exists() bool
	// Clear clears the store
	Clear(ctx context.Context) error
	// GetElementPath determines the path to an element in the store.
	GetElementPath(elementName string) (path string)
	// HasElement States whether elementName is in the store or not.
	HasElement(elementName string) bool
	// GetFilesystem returns the store's filesystem
	GetFilesystem() filesystem.FS
}
