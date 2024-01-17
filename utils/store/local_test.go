/*
 * Copyright (C) 2020-2024 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package store

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ARM-software/golang-utils/utils/filesystem"
)

func TestNewLocalTemporaryStore(t *testing.T) {
	localTestDir, err := filesystem.TempDirInTempDir("test-local-store")
	require.NoError(t, err)
	defer func() { _ = filesystem.Rm(localTestDir) }()
	tests := []struct {
		store IStore
	}{
		{
			store: NewLocalStore(filepath.Join(localTestDir, faker.Word())),
		},
		{
			store: NewLocalTemporaryStore(),
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			store := test.store
			assert.False(t, store.Exists())
			assert.False(t, store.HasElement(faker.Word()))
			assert.Empty(t, store.GetElementPath(faker.Word()))
			require.NoError(t, store.Clear(context.TODO()))
			require.NoError(t, store.Create(context.TODO()))
			assert.True(t, store.Exists())
			element := faker.Word()
			assert.False(t, store.HasElement(element))
			elementPath := store.GetElementPath(element)
			assert.NotEmpty(t, elementPath)
			fullPath := filepath.Clean(filepath.Join(store.GetPath(), elementPath))
			elementPath2 := store.GetElementPath(fullPath)
			assert.NotEmpty(t, elementPath2)
			assert.Equal(t, elementPath, elementPath2)
			// Adding an element in the store
			require.NoError(t, store.GetFilesystem().MkDir(fullPath))
			assert.True(t, store.HasElement(element))
			require.NoError(t, store.Clear(context.TODO()))
			assert.False(t, store.HasElement(element))
			require.NoError(t, store.Close())
		})
	}
}
