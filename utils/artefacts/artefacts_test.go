/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package artefacts

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ARM-software/embedded-development-services-client/client"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/commonerrors/errortest"
	"github.com/ARM-software/golang-utils/utils/field"
	"github.com/ARM-software/golang-utils/utils/filesystem"
	"github.com/ARM-software/golang-utils/utils/hashing"
	"github.com/ARM-software/golang-utils/utils/safeio"
)

type testArtefact struct {
	name             string
	path             string
	embeddedResource bool
}

func newTestArtefact(t *testing.T, tmpDir, artefactContent string, embeddedResource bool) *testArtefact {
	path, err := filesystem.TouchTempFile(tmpDir, "artefact")
	require.NoError(t, err)
	if len(artefactContent) > 0 {
		err = filesystem.WriteFile(path, []byte(artefactContent), 0777)
	} else {
		err = filesystem.GetGlobalFileSystem().Touch(path)
	}
	require.NoError(t, err)

	return &testArtefact{
		name:             filepath.Base(path),
		path:             path,
		embeddedResource: embeddedResource,
	}
}

func (t *testArtefact) testGetArtefactManagers(ctx context.Context, _ string) (a *client.ArtefactManagerCollection, resp *http.Response, err error) {
	item, err := t.fetchTestArtefact(ctx)
	if err != nil {
		return
	}
	if t.embeddedResource {
		a = &client.ArtefactManagerCollection{
			Embedded: &client.EmbeddedArtefactManagerItems{Item: []client.ArtefactManagerItem{*item}},
			Links:    *client.NewNullableHalCollectionLinks(client.NewHalCollectionLinksWithDefaults()),
			Metadata: *client.NewNullablePagingMetadata(&client.PagingMetadata{
				Count:  1,
				Ctime:  time.Now(),
				Etime:  client.NullableTime{},
				Limit:  6,
				Mtime:  time.Now(),
				Offset: 0,
				Total:  1,
			}),
			Name:  "list of artefacts",
			Title: faker.Name(),
		}
	} else {
		links := client.NewHalCollectionLinksWithDefaults()
		links.Item = []client.HalLinkData{client.HalLinkData{
			Href:  fmt.Sprintf("/test/%v", item.Name),
			Name:  field.ToOptionalString(item.Name),
			Title: field.ToOptionalString(faker.Name()),
		}}
		a = &client.ArtefactManagerCollection{
			Embedded: nil,
			Links:    *client.NewNullableHalCollectionLinks(links),
			Metadata: *client.NewNullablePagingMetadata(&client.PagingMetadata{
				Count:  1,
				Ctime:  time.Now(),
				Etime:  client.NullableTime{},
				Limit:  6,
				Mtime:  time.Now(),
				Offset: 0,
				Total:  1,
			}),
			Name:  "list of artefacts",
			Title: faker.Name(),
		}
	}
	resp = &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(safeio.NewByteReader(ctx, []byte("hello")))}
	return
}

func (t *testArtefact) testGetArtefactManager(ctx context.Context, _, artefact string) (a *client.ArtefactManagerItem, resp *http.Response, err error) {
	if artefact == t.name {
		a, err = t.fetchTestArtefact(ctx)
		resp = &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(safeio.NewByteReader(ctx, []byte("hello")))}
	} else {
		err = commonerrors.ErrNotFound
	}
	return
}

func (t *testArtefact) fetchTestArtefact(ctx context.Context) (a *client.ArtefactManagerItem, err error) {
	fileHasher, subErr := filesystem.NewFileHash(hashing.HashSha256)
	if subErr != nil {
		err = subErr
		return
	}

	f, subErr := filesystem.OpenFile(t.path, os.O_RDWR, 0777)
	if subErr != nil {
		err = subErr
		return
	}

	hash, subErr := fileHasher.CalculateWithContext(ctx, f)
	if subErr != nil {
		err = subErr
		return
	}

	size, subErr := filesystem.GetFileSize(t.path)
	if subErr != nil {
		err = subErr
		return
	}

	a = &client.ArtefactManagerItem{
		Name:  t.name,
		Title: *client.NewNullableString(field.ToOptionalString(t.name)),
		Hash:  *client.NewNullableString(&hash),
		Size:  &size,
	}
	return
}

func (t *testArtefact) testGetOutputArtefact(ctx context.Context, _, artefact string) (f *os.File, resp *http.Response, err error) {
	if artefact == t.name {
		f, err = os.Open(t.path)
		if err != nil {
			return
		}

		return f, &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(safeio.NewByteReader(ctx, []byte("hello")))}, nil
	}

	return nil, &http.Response{StatusCode: http.StatusNotFound, Body: io.NopCloser(safeio.NewByteReader(ctx, []byte("hello")))}, commonerrors.ErrNotFound
}

func newTestArtefactManager(t *testing.T, tmpDir, artefactContent string, linksOnly bool) (IArtefactManager, *testArtefact) {
	testArtefact := newTestArtefact(t, tmpDir, artefactContent, !linksOnly)
	return NewArtefactManager(testArtefact.testGetArtefactManagers, nil, testArtefact.testGetArtefactManager, testArtefact.testGetOutputArtefact), testArtefact
}

func newTestArtefactManagerWithEmbeddedResources(t *testing.T, tmpDir, artefactContent string) (IArtefactManager, *testArtefact) {
	return newTestArtefactManager(t, tmpDir, artefactContent, false)
}

func TestDetermineDestination(t *testing.T) {
	outputDir := strings.ReplaceAll(faker.Sentence(), " ", "//") + "                  "
	cleanedOutputDir := filepath.Clean(outputDir)

	tests := []struct {
		item             client.ArtefactManagerItem
		maintainTree     bool
		outputDir        string
		expectedFileName string
		expectedDir      string
	}{
		{
			item: client.ArtefactManagerItem{
				ExtraMetadata: nil,
				Name:          faker.Name(),
				Size:          nil,
				Title:         *client.NewNullableString(field.ToOptionalString("test.j")),
			},
			maintainTree:     false,
			outputDir:        outputDir,
			expectedFileName: "test.j",
			expectedDir:      cleanedOutputDir,
		},
		{
			item: client.ArtefactManagerItem{
				ExtraMetadata: nil,
				Name:          faker.Name(),
				Size:          nil,
				Title:         *client.NewNullableString(field.ToOptionalString("test.j")),
			},
			maintainTree:     true,
			outputDir:        outputDir,
			expectedFileName: "test.j",
			expectedDir:      cleanedOutputDir,
		},
		{
			item: client.ArtefactManagerItem{
				ExtraMetadata: nil,
				Name:          faker.Name(),
				Size:          nil,
				Title:         *client.NewNullableString(field.ToOptionalString("cool+blog&about%2Cstuff.yep")),
			},
			maintainTree:     true,
			outputDir:        outputDir,
			expectedFileName: "cool+blog&about,stuff.yep",
			expectedDir:      cleanedOutputDir,
		},
		{
			item: client.ArtefactManagerItem{
				ExtraMetadata: &map[string]string{faker.Name(): faker.Sentence()},
				Name:          faker.Name(),
				Size:          nil,
				Title:         *client.NewNullableString(field.ToOptionalString("cool+blog&about%2Cstuff.yep")),
			},
			maintainTree:     true,
			outputDir:        outputDir,
			expectedFileName: "cool+blog&about,stuff.yep",
			expectedDir:      cleanedOutputDir,
		},
		{
			item: client.ArtefactManagerItem{
				ExtraMetadata: &map[string]string{relativePathKey: "        test/1     "},
				Name:          faker.Name(),
				Size:          nil,
				Title:         *client.NewNullableString(field.ToOptionalString("cool+blog&about%2Cstuff.yep")),
			},
			maintainTree:     true,
			outputDir:        outputDir,
			expectedFileName: "cool+blog&about,stuff.yep",
			expectedDir:      filepath.Join(cleanedOutputDir, "test", "1"),
		},
		{
			item: client.ArtefactManagerItem{
				ExtraMetadata: &map[string]string{relativePathKey: "        test/1/cool+blog&about,stuff.yep     "},
				Name:          faker.Name(),
				Size:          nil,
				Title:         *client.NewNullableString(field.ToOptionalString("cool+blog&about%2Cstuff.yep")),
			},
			maintainTree:     true,
			outputDir:        outputDir,
			expectedFileName: "cool+blog&about,stuff.yep",
			expectedDir:      filepath.Join(cleanedOutputDir, "test", "1"),
		},
		{
			item: client.ArtefactManagerItem{
				ExtraMetadata: &map[string]string{relativePathKey: "        test/1/cool+blog&about%2Cstuff.yep     "},
				Name:          faker.Name(),
				Size:          nil,
				Title:         *client.NewNullableString(field.ToOptionalString("cool+blog&about%2Cstuff.yep")),
			},
			maintainTree:     true,
			outputDir:        outputDir,
			expectedFileName: "cool+blog&about,stuff.yep",
			expectedDir:      filepath.Join(cleanedOutputDir, "test", "1"),
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run(fmt.Sprintf("%d_%s", i, test.expectedFileName), func(t *testing.T) {
			fileName, fileDest, err := determineArtefactDestination(test.outputDir, test.maintainTree, &test.item)
			require.NoError(t, err)
			assert.Equal(t, test.expectedFileName, fileName)
			assert.Equal(t, test.expectedDir, fileDest)
		})
	}
}

func TestArtefactDownload(t *testing.T) {
	t.Run("Happy download artefact", func(t *testing.T) {
		tmpDir, err := filesystem.TempDirInTempDir("test-artefact-")
		require.NoError(t, err)
		defer func() { _ = filesystem.Rm(tmpDir) }()
		m, a := newTestArtefactManagerWithEmbeddedResources(t, tmpDir, faker.Sentence())

		out := t.TempDir()
		err = m.DownloadJobArtefactFromLink(context.Background(), faker.Word(), out, &client.HalLinkData{
			Name: &a.name,
		})
		require.NoError(t, err)

		require.FileExists(t, filepath.Join(out, a.name))
		expectedContents, err := filesystem.ReadFile(a.path)
		require.NoError(t, err)
		actualContents, err := filesystem.ReadFile(filepath.Join(out, a.name))
		require.NoError(t, err)
		assert.Equal(t, expectedContents, actualContents)
	})
	t.Run("Happy download artefact and keep tree", func(t *testing.T) {
		tmpDir, err := filesystem.TempDirInTempDir("test-artefact-with-tree-")
		require.NoError(t, err)
		defer func() { _ = filesystem.Rm(tmpDir) }()
		m, a := newTestArtefactManagerWithEmbeddedResources(t, tmpDir, faker.Sentence())

		out := t.TempDir()
		err = m.DownloadJobArtefactFromLinkWithTree(context.Background(), faker.Word(), true, out, &client.HalLinkData{
			Name: &a.name,
		})
		require.NoError(t, err)

		require.FileExists(t, filepath.Join(out, a.name))
		expectedContents, err := filesystem.ReadFile(a.path)
		require.NoError(t, err)
		actualContents, err := filesystem.ReadFile(filepath.Join(out, a.name))
		require.NoError(t, err)
		assert.Equal(t, expectedContents, actualContents)
	})
	t.Run("Happy list artefacts links", func(t *testing.T) {
		tmpDir, err := filesystem.TempDirInTempDir("test-artefact-")
		require.NoError(t, err)
		defer func() { _ = filesystem.Rm(tmpDir) }()
		m, _ := newTestArtefactManager(t, tmpDir, faker.Sentence(), true)

		iter, err := m.ListJobArtefacts(context.Background(), faker.Word())
		require.NoError(t, err)
		assert.True(t, iter.HasNext())
		element, err := iter.GetNext()
		require.NoError(t, err)
		assert.NotNil(t, element)
		_, ok := element.(*client.HalLinkData)
		assert.True(t, ok)
		assert.False(t, iter.HasNext())
	})

	t.Run("Happy list artefacts", func(t *testing.T) {
		tmpDir, err := filesystem.TempDirInTempDir("test-artefact-")
		require.NoError(t, err)
		defer func() { _ = filesystem.Rm(tmpDir) }()
		m, _ := newTestArtefactManagerWithEmbeddedResources(t, tmpDir, faker.Sentence())

		iter, err := m.ListJobArtefacts(context.Background(), faker.Word())
		require.NoError(t, err)
		assert.True(t, iter.HasNext())
		element, err := iter.GetNext()
		require.NoError(t, err)
		assert.NotNil(t, element)
		_, ok := element.(*client.ArtefactManagerItem)
		assert.True(t, ok)
		assert.False(t, iter.HasNext())
	})

	t.Run("Happy download all artefacts", func(t *testing.T) {
		tmpDir, err := filesystem.TempDirInTempDir("test-artefact-")
		require.NoError(t, err)
		defer func() { _ = filesystem.Rm(tmpDir) }()
		m, a := newTestArtefactManager(t, tmpDir, faker.Sentence(), true)

		out := t.TempDir()
		err = m.DownloadAllJobArtefacts(context.Background(), faker.Word(), out)
		require.NoError(t, err)

		require.FileExists(t, filepath.Join(out, a.name))
		expectedContents, err := filesystem.ReadFile(a.path)
		require.NoError(t, err)
		actualContents, err := filesystem.ReadFile(filepath.Join(out, a.name))
		require.NoError(t, err)
		assert.Equal(t, expectedContents, actualContents)
	})

	t.Run("Invalid Artefact", func(t *testing.T) {
		tmpDir, err := filesystem.TempDirInTempDir("test-artefact-")
		require.NoError(t, err)
		defer func() { _ = filesystem.Rm(tmpDir) }()
		m, _ := newTestArtefactManagerWithEmbeddedResources(t, tmpDir, faker.Sentence())

		out := t.TempDir()
		err = m.DownloadJobArtefactFromLink(context.Background(), faker.Word(), out, &client.HalLinkData{
			Name: field.ToOptionalString(faker.Word()),
		})
		errortest.AssertErrorDescription(t, err, "cannot fetch artefact's manager")
	})

	t.Run("Empty Artefact", func(t *testing.T) {
		tmpDir, err := filesystem.TempDirInTempDir("test-artefact-")
		require.NoError(t, err)
		defer func() { _ = filesystem.Rm(tmpDir) }()
		m, a := newTestArtefactManagerWithEmbeddedResources(t, tmpDir, "")

		out := t.TempDir()
		err = m.DownloadJobArtefactFromLink(context.Background(), faker.Word(), out, &client.HalLinkData{
			Name: &a.name,
		})
		errortest.AssertError(t, err, commonerrors.ErrEmpty)
	})
}
