/*
 * Copyright (C) 2020-2024 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package artefacts

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ARM-software/embedded-development-services-client/client"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/commonerrors/errortest"
	"github.com/ARM-software/golang-utils/utils/field"
	"github.com/ARM-software/golang-utils/utils/filesystem"
	"github.com/ARM-software/golang-utils/utils/hashing"
)

type testArtefact struct {
	name string
	path string
}

func newTestArtefact(t *testing.T, artefactContent string) *testArtefact {
	tmpDir := t.TempDir()

	path, err := filesystem.TouchTempFile(tmpDir, "artefact")
	require.NoError(t, err)

	f, err := filesystem.OpenFile(path, os.O_RDWR, 0777)
	require.NoError(t, err)

	_, err = f.Write([]byte(artefactContent))
	require.NoError(t, err)

	of := filesystem.ConvertToOSFile(f)
	return &testArtefact{
		name: filepath.Base(of.Name()),
		path: path,
	}
}

func (t *testArtefact) testGetArtefactManager(ctx context.Context, _, artefact string) (a *client.ArtefactManagerItem, resp *http.Response, err error) {
	if artefact == t.name {
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

		return &client.ArtefactManagerItem{
			Name:  t.name,
			Title: *client.NewNullableString(field.ToOptionalString(t.name)),
			Hash:  *client.NewNullableString(&hash),
			Size:  &size,
		}, &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader([]byte("hello")))}, nil
	}
	return nil, nil, nil
}

func (t *testArtefact) testGetOutputArtefact(_ context.Context, _, artefact string) (f *os.File, resp *http.Response, err error) {
	if artefact == t.name {
		f, err = os.Open(t.path)
		if err != nil {
			return
		}

		return f, &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader([]byte("hello")))}, nil
	}

	return nil, &http.Response{StatusCode: http.StatusNotFound, Body: io.NopCloser(bytes.NewReader([]byte("hello")))}, commonerrors.ErrNotFound
}

func newTestArtefactManager(t *testing.T, artefactContent string) (*ArtefactManager, *testArtefact) {
	testArtefact := newTestArtefact(t, artefactContent)
	return NewArtefactManager(testArtefact.testGetArtefactManager, testArtefact.testGetOutputArtefact), testArtefact
}

func TestArtefactDownload(t *testing.T) {
	t.Run("Happy", func(t *testing.T) {
		m, a := newTestArtefactManager(t, faker.Sentence())

		out := t.TempDir()
		err := m.DownloadJobArtefact(context.Background(), faker.Word(), out, client.HalLinkData{
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

	t.Run("Invalid Artefact", func(t *testing.T) {
		m, _ := newTestArtefactManager(t, faker.Sentence())

		out := t.TempDir()
		err := m.DownloadJobArtefact(context.Background(), faker.Word(), out, client.HalLinkData{
			Name: field.ToOptionalString(faker.Word()),
		})
		errortest.AssertErrorDescription(t, err, "cannot fetch artefact's manager")
	})

	t.Run("Empty Artefact", func(t *testing.T) {
		m, a := newTestArtefactManager(t, "")

		out := t.TempDir()
		err := m.DownloadJobArtefact(context.Background(), faker.Word(), out, client.HalLinkData{
			Name: &a.name,
		})
		errortest.AssertError(t, err, commonerrors.ErrEmpty)
	})
}
