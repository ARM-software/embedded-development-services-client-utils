/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package artefacts

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/api"
	paginationUtils "github.com/ARM-software/embedded-development-services-client-utils/utils/pagination"
	"github.com/ARM-software/embedded-development-services-client/client"
	"github.com/ARM-software/golang-utils/utils/collection/pagination"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/filesystem"
	"github.com/ARM-software/golang-utils/utils/hashing"
	"github.com/ARM-software/golang-utils/utils/parallelisation"
	"github.com/ARM-software/golang-utils/utils/reflection"
	"github.com/ARM-software/golang-utils/utils/safeio"
)

const relativePathKey = "Relative Path"

type (
	// GetArtefactManagersFirstPageFunc defines the function which can retrieve the first page of artefact managers.
	GetArtefactManagersFirstPageFunc[D ILinkData, L ILinks[D], C ICollection[D, L]] = func(ctx context.Context, job string) (C, *http.Response, error)
	// FollowLinkToArtefactManagersPageFunc is a function able to follow a link to an artefact manager page.
	FollowLinkToArtefactManagersPageFunc[D ILinkData, L ILinks[D], C ICollection[D, L]] = func(ctx context.Context, link D) (C, *http.Response, error)
	// GetArtefactManagerFunc is a function which retrieves information about an artefact manager.
	GetArtefactManagerFunc[M IManager] = func(ctx context.Context, job, artefact string) (M, *http.Response, error)
	// GetArtefactContentFunc is a function able to return the content of any artefact managers.
	GetArtefactContentFunc = func(ctx context.Context, job, artefactID string) (*os.File, *http.Response, error)
)

func determineArtefactDestination[M IManager](outputDir string, maintainTree bool, item M) (artefactFileName string, destinationDir string, err error) {
	if any(item) == nil {
		err = commonerrors.UndefinedVariable("artefact item")
		return
	}
	artefactManagerName := item.GetName()
	if artefactManagerName == "" {
		err = commonerrors.UndefinedVariable("artefact name")
		return
	}
	rawFileName := artefactManagerName
	if item.HasTitle() {
		rawFileName = item.GetTitle()
	}
	artefactFileName = rawFileName
	if unescapedName, err := url.PathUnescape(rawFileName); err == nil {
		artefactFileName = unescapedName
	}
	destinationDir = filepath.Clean(outputDir)
	if !maintainTree {
		return
	}

	if item.HasExtraMetadata() {
		m := item.GetExtraMetadata()
		treePath, ok := m[relativePathKey]
		if !ok {
			return
		}
		treePath = strings.TrimSpace(treePath)
		if strings.HasSuffix(treePath, rawFileName) || strings.HasSuffix(treePath, artefactFileName) {
			treePath = filepath.Dir(treePath)
		}
		destinationDir = filepath.Clean(filepath.Join(outputDir, treePath))
	}
	return
}

type ArtefactManager[
	M IManager,
	D ILinkData,
	L ILinks[D],
	C ICollection[D, L],
] struct {
	getArtefactManagerFunc            GetArtefactManagerFunc[M]
	getArtefactContentFunc            GetArtefactContentFunc
	getArtefactManagersFirstPageFunc  GetArtefactManagersFirstPageFunc[D, L, C]
	getArtefactManagersFollowLinkFunc FollowLinkToArtefactManagersPageFunc[D, L, C]
}

// NewArtefactManager returns an artefact manager.
func NewArtefactManager[
	M IManager,
	D ILinkData,
	L ILinks[D],
	C ICollection[D, L],
](
	getArtefactManagersFirstPage GetArtefactManagersFirstPageFunc[D, L, C],
	getArtefactsManagersPage FollowLinkToArtefactManagersPageFunc[D, L, C],
	getArtefactManager GetArtefactManagerFunc[M],
	getOutputArtefact GetArtefactContentFunc) IArtefactManager[M, D] {
	return &ArtefactManager[M, D, L, C]{
		getArtefactManagerFunc:            getArtefactManager,
		getArtefactContentFunc:            getOutputArtefact,
		getArtefactManagersFirstPageFunc:  getArtefactManagersFirstPage,
		getArtefactManagersFollowLinkFunc: getArtefactsManagersPage,
	}
}
func (m *ArtefactManager[M, D, L, C]) DownloadJobArtefact(ctx context.Context, jobName string, outputDirectory string, artefactManager M) (err error) {
	return m.DownloadJobArtefactWithTree(ctx, jobName, false, outputDirectory, artefactManager)
}

func (m *ArtefactManager[M, D, L, C]) DownloadJobArtefactWithTree(ctx context.Context, jobName string, maintainTreeLocation bool, outputDirectory string, artefactManager M) (err error) {
	err = parallelisation.DetermineContextError(ctx)
	if err != nil {
		return
	}
	if m.getArtefactManagerFunc == nil || m.getArtefactContentFunc == nil {
		err = commonerrors.New(commonerrors.ErrUndefined, "function to retrieve an artefact manager was not properly defined")
		return
	}

	err = filesystem.MkDir(outputDirectory)
	if err != nil {
		err = commonerrors.WrapErrorf(commonerrors.ErrUnexpected, err, "failed creating the output directory [%v] for job artefact", outputDirectory)
		return
	}

	fileHasher, err := filesystem.NewFileHash(hashing.HashSha256)
	if err != nil {
		return
	}
	if any(artefactManager) == nil {
		err = commonerrors.UndefinedVariable("artefact manager")
		return
	}

	artefactManagerName := artefactManager.GetName()
	if artefactManagerName == "" {
		err = commonerrors.UndefinedVariable("artefact name")
		return
	}

	expectedSizePtr, ok := artefactManager.GetSizeOk()
	if !ok {
		err = commonerrors.Newf(commonerrors.ErrUndefined, "could not fetch artefact's size from artefact's manager [%v]", artefactManagerName)
		return
	}
	expectedSize := *expectedSizePtr

	expectedHashPtr, ok := artefactManager.GetHashOk()
	if !ok {
		err = commonerrors.Newf(commonerrors.ErrUndefined, "could not fetch artefact's hash from artefact's manager [%v]", artefactManagerName)
		return
	}
	expectedHash := *expectedHashPtr

	artefactFilename, artefactDestDir, err := determineArtefactDestination(outputDirectory, maintainTreeLocation, artefactManager)
	if err != nil {
		return
	}
	if reflection.IsEmpty(artefactFilename) {
		err = commonerrors.UndefinedVariable("artefact filename")
		return
	}
	err = filesystem.MkDir(artefactDestDir)
	if err != nil {
		err = commonerrors.WrapErrorf(commonerrors.ErrUnexpected, err, "failed creating the output directory [%v] for job artefact", artefactDestDir)
		return
	}
	artefact, err := api.CallAndCheckSuccess[os.File](ctx, fmt.Sprintf("cannot fetch generated artefact [%v]", artefactFilename), func(fCtx context.Context) (*os.File, *http.Response, error) {
		return m.getArtefactContentFunc(fCtx, jobName, artefactManagerName)
	})
	defer func() {
		if artefact != nil {
			_ = artefact.Close()
		}
	}()
	if err != nil {
		return
	}
	destination, err := filesystem.CreateFile(filepath.Join(artefactDestDir, artefactFilename))
	if err != nil {
		err = commonerrors.WrapErrorf(commonerrors.ErrUnexpected, err, "could not create a location to store generated artefact [%v]", artefactFilename)
		return
	}
	defer func() { _ = destination.Close() }()

	actualSize, err := safeio.CopyDataWithContext(ctx, artefact, destination)
	if err != nil {
		err = commonerrors.WrapErrorf(commonerrors.ErrUnexpected, err, "failed to copy artefact [%v]", artefactFilename)
		return
	}
	if actualSize == 0 {
		err = commonerrors.Newf(commonerrors.ErrEmpty, "problem with artefact [%v]", artefactFilename)
		return
	}
	if actualSize != expectedSize {
		err = commonerrors.Newf(commonerrors.ErrCondition, "artefact [%v] size '%v' does not match expected '%v'", artefactFilename, actualSize, expectedSize)
		return
	}

	// reset offset for hashing entire contents
	_, err = destination.Seek(0, 0)
	if err != nil {
		err = commonerrors.WrapError(commonerrors.ErrUnexpected, err, "could not reset destination file")
		return
	}

	actualHash, err := fileHasher.CalculateWithContext(ctx, destination)
	if err != nil {
		err = commonerrors.WrapError(commonerrors.ErrUnexpected, err, "could not calculate hash of destination file")
	}
	if actualHash != expectedHash {
		err = commonerrors.Newf(commonerrors.ErrCondition, "artefact [%v] hash '%v' does not match expected '%v'", artefactFilename, actualHash, expectedHash)
		return
	}

	err = parallelisation.DetermineContextError(ctx)
	return

}
func (m *ArtefactManager[M, D, L, C]) DownloadJobArtefactFromLink(ctx context.Context, jobName string, outputDirectory string, artefactManagerItemLink D) error {
	return m.DownloadJobArtefactFromLinkWithTree(ctx, jobName, false, outputDirectory, artefactManagerItemLink)
}

func (m *ArtefactManager[M, D, L, C]) DownloadJobArtefactFromLinkWithTree(ctx context.Context, jobName string, maintainTreeLocation bool, outputDirectory string, artefactManagerItemLink D) (err error) {
	err = parallelisation.DetermineContextError(ctx)
	if err != nil {
		return
	}
	if m.getArtefactManagerFunc == nil || m.getArtefactContentFunc == nil {
		err = commonerrors.New(commonerrors.ErrUndefined, "function to retrieve an artefact manager was not properly defined")
		return
	}
	if any(artefactManagerItemLink) == nil {
		err = commonerrors.UndefinedVariable("artefact link")
		return
	}

	artefactManagerName := artefactManagerItemLink.GetName()
	if reflection.IsEmpty(artefactManagerName) {
		err = commonerrors.UndefinedVariable("artefact name")
		return
	}
	artefactManager, err := api.GenericCallAndCheckSuccess[M](ctx, fmt.Sprintf("cannot fetch artefact's manager [%v]", artefactManagerName), func(fCtx context.Context) (M, *http.Response, error) {
		return m.getArtefactManagerFunc(fCtx, jobName, artefactManagerName)
	})
	if err != nil {
		return
	}
	err = m.DownloadJobArtefactWithTree(ctx, jobName, maintainTreeLocation, outputDirectory, artefactManager)
	return
}

func (m *ArtefactManager[M, D, L, C]) ListJobArtefacts(ctx context.Context, jobName string) (pagination.IPaginatorAndPageFetcher, error) {
	err := parallelisation.DetermineContextError(ctx)
	if err != nil {
		return nil, err
	}
	return pagination.NewStaticPagePaginator(ctx, func(context.Context) (pagination.IStaticPage, error) {
		return m.fetchJobArtefactsFirstPage(ctx, jobName)
	}, m.fetchJobArtefactsNextPage)
}

func (m *ArtefactManager[M, D, L, C]) fetchJobArtefactsFirstPage(ctx context.Context, jobName string) (page pagination.IStaticPage, err error) {
	if m.getArtefactManagersFirstPageFunc == nil {
		err = commonerrors.New(commonerrors.ErrUndefined, "function to retrieve artefact managers was not properly defined")
		return
	}
	if reflection.IsEmpty(jobName) {
		err = commonerrors.UndefinedVariable("job identifier")
		return
	}
	clientPage, err := api.GenericCallAndCheckSuccess[client.IStaticPage](ctx, fmt.Sprintf("could not list artefact managers for job [%v]", jobName), func(fCtx context.Context) (client.IStaticPage, *http.Response, error) {
		return m.getArtefactManagersFirstPageFunc(fCtx, jobName)
	})
	if err == nil {
		page = paginationUtils.ToPage(clientPage)
	}
	return
}

func (m *ArtefactManager[M, D, L, C]) fetchJobArtefactsNextPage(ctx context.Context, currentPage pagination.IStaticPage) (nextPage pagination.IStaticPage, err error) {
	err = parallelisation.DetermineContextError(ctx)
	if err != nil {
		return
	}
	if currentPage == nil {
		return
	}
	if m.getArtefactManagersFollowLinkFunc == nil {
		err = commonerrors.New(commonerrors.ErrUndefined, "function to retrieve artefact managers was not properly defined")
		return
	}
	unwrappedPage := paginationUtils.UnwrapPage(currentPage)
	if unwrappedPage == nil {
		err = commonerrors.New(commonerrors.ErrUnexpected, "returned artefact managers page is empty")
		return
	}
	page, ok := unwrappedPage.(C)
	if !ok {
		err = commonerrors.Newf(commonerrors.ErrUnexpected, " returned artefact managers page[%T] is not of the expected type [%v]", currentPage, "*ArtefactManagerCollection")
		return
	}
	links, has := page.GetLinksOk()
	if !has {
		err = commonerrors.New(commonerrors.ErrUnexpected, "returned page of artefact managers has no links")
		return
	}
	if !links.HasNext() {
		err = commonerrors.New(commonerrors.ErrUnexpected, "returned page of artefact managers has no `next` link")
		return
	}
	link := links.GetNextP()
	if reflection.IsEmpty(link) {
		err = commonerrors.UndefinedVariable("`next` link")
		return
	}
	clientPage, err := api.GenericCallAndCheckSuccess[client.IStaticPage](ctx, fmt.Sprintf("could not follow `next` link [%v]", link), func(fCtx context.Context) (client.IStaticPage, *http.Response, error) {
		return m.getArtefactManagersFollowLinkFunc(fCtx, link)
	})
	if err == nil {
		nextPage = paginationUtils.ToPage(clientPage)
	}
	return
}

func (m *ArtefactManager[M, D, L, C]) DownloadAllJobArtefacts(ctx context.Context, jobName string, outputDirectory string) error {
	return m.DownloadAllJobArtefactsWithTree(ctx, jobName, false, outputDirectory)
}

func (m *ArtefactManager[M, D, L, C]) DownloadAllJobArtefactsWithTree(ctx context.Context, jobName string, maintainTreeStructure bool, outputDirectory string) (err error) {
	return m.DownloadAllJobArtefactsWithOptions(ctx, jobName, outputDirectory, WithMaintainStructure(maintainTreeStructure))
}

func (m *ArtefactManager[M, D, L, C]) DownloadAllJobArtefactsWithOptions(ctx context.Context, jobName string, outputDirectory string, opts ...DownloadOption) (err error) {
	err = parallelisation.DetermineContextError(ctx)
	if err != nil {
		return
	}

	dlOpts := NewDownloadOptions(opts...)
	err = filesystem.MkDir(outputDirectory)
	if err != nil {
		err = commonerrors.WrapErrorf(commonerrors.ErrUnexpected, err, "failed creating the output directory [%v] for job artefacts", outputDirectory)
		return
	}
	paginator, err := m.ListJobArtefacts(ctx, jobName)
	if err != nil {
		return
	}
	stop := paginator.Stop()
	defer stop()
	var collatedDownloadErrors []error
	for {
		if !paginator.HasNext() {
			if len(collatedDownloadErrors) > 0 {
				err = commonerrors.Join(collatedDownloadErrors...)
			}
			return
		}
		item, subErr := paginator.GetNext()
		if subErr != nil {
			err = commonerrors.WrapError(commonerrors.ErrUnexpected, subErr, "failed getting information about job artefacts")
			return
		}
		var artefactName string
		var downloadErr error
		artefactLink, ok := item.(D)
		if ok {
			artefactName = artefactLink.GetName()
			downloadErr = m.DownloadJobArtefactFromLinkWithTree(ctx, jobName, dlOpts.MaintainTreeStructure, outputDirectory, artefactLink)
		} else {
			artefactManager, isManager := item.(M)
			if isManager {
				artefactName = artefactManager.GetName()
				downloadErr = m.DownloadJobArtefactWithTree(ctx, jobName, dlOpts.MaintainTreeStructure, outputDirectory, artefactManager)
			} else {
				downloadErr = commonerrors.New(commonerrors.ErrMarshalling, "the type of the response from service cannot be interpreted")
			}
		}

		if downloadErr != nil {
			if dlOpts.StopOnFirstError {
				err = downloadErr
				return
			}
			collatedDownloadErrors = append(collatedDownloadErrors, downloadErr)
			if dlOpts.Logger != nil {
				dlOpts.Logger.LogError(downloadErr)
			}
		} else if !reflection.IsEmpty(artefactName) {
			if dlOpts.Logger != nil {
				dlOpts.Logger.Log(fmt.Sprintf("downloading %s", artefactName))
			}
		}
	}
}
