/*
 * Copyright (C) 2020-2024 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package artefacts

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	utils "github.com/ARM-software/embedded-development-services-client-utils/utils/api"
	"github.com/ARM-software/embedded-development-services-client/client"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/filesystem"
	"github.com/ARM-software/golang-utils/utils/hashing"
	"github.com/ARM-software/golang-utils/utils/parallelisation"
)

type (
	getArtefactManagerFunc = func(ctx context.Context, job, artefact string) (*client.ArtefactManagerItem, *http.Response, error)
	getOutputArtefactFunc  = func(ctx context.Context, job, artefactID string) (*os.File, *http.Response, error)
)

type ArtefactManager struct {
	getArtefactManager getArtefactManagerFunc
	getOutputArtefact  getOutputArtefactFunc
}

func NewArtefactManager(getArtefactManager getArtefactManagerFunc, getOutputArtefact getOutputArtefactFunc) *ArtefactManager {
	return &ArtefactManager{
		getArtefactManager: getArtefactManager,
		getOutputArtefact:  getOutputArtefact,
	}
}

func (m *ArtefactManager) DownloadJobArtefact(ctx context.Context, jobName string, outputDirectory string, artefactManagerItem client.HalLinkData) (err error) {
	artefactManagerName := artefactManagerItem.GetName()
	artefactManager, resp, apierr := m.getArtefactManager(ctx, jobName, artefactManagerName)
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()
	err = utils.CheckAPICallSuccess(ctx, fmt.Sprintf("cannot fetch artefact manager [%v]", artefactManager), resp, apierr)
	if err != nil {
		return
	}

	artefactFilenamePtr, ok := artefactManager.GetTitleOk()
	if !ok {
		err = fmt.Errorf("%w: could not fetch artefact title from artefact manager [%v]", commonerrors.ErrUndefined, artefactManagerName)
		return
	}
	artefactFilename := *artefactFilenamePtr

	expectedSizePtr, ok := artefactManager.GetSizeOk()
	if !ok {
		err = fmt.Errorf("%w: could not fetch artefact size from artefact manager [%v]", commonerrors.ErrUndefined, artefactManagerName)
		return
	}
	expectedSize := *expectedSizePtr

	expectedHashPtr, ok := artefactManager.GetHashOk()
	if !ok {
		err = fmt.Errorf("%w: could not fetch artefact hash from artefact manager [%v]", commonerrors.ErrUndefined, artefactManagerName)
		return
	}
	expectedHash := *expectedHashPtr

	artefactIDPtr, ok := artefactManager.GetNameOk()
	if !ok {
		err = fmt.Errorf("%w: could not fetch artefact name from artefact manager [%v]", commonerrors.ErrUndefined, artefactManagerName)
		return
	}
	artefactID := *artefactIDPtr

	artefact, resp, apierr := m.getOutputArtefact(ctx, jobName, artefactID)
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
		if artefact != nil {
			_ = artefact.Close()
		}
	}()

	err = utils.CheckAPICallSuccess(ctx, fmt.Sprintf("cannot fetch generated artefact [%v]", artefactFilename), resp, apierr)
	if err != nil {
		return
	}

	unescapedName, err := url.PathUnescape(artefactFilename)
	if err == nil {
		artefactFilename = unescapedName
	}

	destinationPath := filepath.Join(outputDirectory, artefactFilename)
	destination, err := filesystem.CreateFile(destinationPath)
	if err != nil {
		err = fmt.Errorf("%w: could not create a location to store generated artefact [%v]: %v", commonerrors.ErrUnexpected, artefactFilename, err.Error())
		return
	}
	defer func() { _ = destination.Close() }()

	fileHasher, err := filesystem.NewFileHash(hashing.HashSha256)
	if err != nil {
		return
	}

	actualSize, err := io.Copy(destination, artefact)
	if err != nil {
		err = fmt.Errorf("%w: failed to copy artefact [%v]: %v", commonerrors.ErrUnexpected, artefactFilename, err.Error())
		return
	}
	if actualSize == 0 {
		err = fmt.Errorf("%w: problem with artefact [%v]", commonerrors.ErrEmpty, artefactFilename)
		return
	}
	if actualSize != expectedSize {
		err = fmt.Errorf("%w: artefact [%v] size '%v' does not match expected '%v'", commonerrors.ErrCondition, artefactFilename, actualSize, expectedSize)
		return
	}

	// reset offset for hashing entire contents
	_, err = destination.Seek(0, 0)
	if err != nil {
		err = fmt.Errorf("%w: could not reset destination file: %v", commonerrors.ErrUnexpected, err.Error())
		return
	}

	actualHash, err := fileHasher.CalculateWithContext(ctx, destination)
	if err != nil {
		err = fmt.Errorf("%w: could not calculate hash of destination file: %v", commonerrors.ErrUnexpected, err.Error())
	}
	if actualHash != expectedHash {
		err = fmt.Errorf("%w: artefact [%v] hash '%v' does not match expected '%v'", commonerrors.ErrCondition, artefactFilename, actualHash, expectedHash)
		return
	}

	err = parallelisation.DetermineContextError(ctx)
	return
}
