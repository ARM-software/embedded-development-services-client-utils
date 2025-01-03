/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package cache

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ServiceCache(t *testing.T) {
	cache := NewServiceCache()
	assert.Empty(t, cache.GetKey())
	key := faker.UUIDHyphenated()
	require.NoError(t, cache.SetKey(key))
	assert.Equal(t, key, cache.GetKey())
	require.NoError(t, cache.Invalidate(context.TODO()))
	assert.Empty(t, cache.GetKey())
	require.NoError(t, cache.SetCacheControl(NoStore))
	require.NoError(t, cache.SetKey(key))
	assert.Empty(t, cache.GetKey())
}
