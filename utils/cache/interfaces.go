/*
 * Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// Package cache defines utilities for caching data or state.
package cache

import "context"

// Mocks are generated using `go generate ./...`
// Add interfaces to the following command for a mock to be generated
//go:generate mockgen -destination=../mocks/mock_$GOPACKAGE.go -package=mocks github.com/ARM-software/embedded-development-services-client-utils/utils/$GOPACKAGE IServerCache

type CacheControl int

const (
	// Apply caching wherever possible
	Apply = iota
	// NoCache directive does not prevent the storing of data but instead prevents the reuse of data without revalidation
	NoCache
	// NoStore directive ensures there is no caching performed at all
	NoStore
)

// IServerCache defines a caching mechanism server-side.
type IServerCache interface {
	// SetCacheControl specifies the caching behaviour.
	SetCacheControl(control CacheControl) error
	// GetCacheControl returns the caching behaviour followed.
	GetCacheControl() CacheControl
	// SetKey sets  an explicit key for restoring and saving the cache
	SetKey(key string) error
	// GetKey returns  an explicit key for restoring a cache
	GetKey() string
	// Invalidate invalidates all cache entries.
	Invalidate(ctx context.Context) error
}
