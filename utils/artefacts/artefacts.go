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
	"github.com/ARM-software/golang-utils/utils/collection/pagination"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/filesystem"
	"github.com/ARM-software/golang-utils/utils/hashing"
	"github.com/ARM-software/golang-utils/utils/parallelisation"
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
		err = fmt.Errorf("%w: missing artefact item", commonerrors.ErrUndefined)
		return
	}
	artefactManagerName := item.GetName()
	if artefactManagerName == "" {
		err = fmt.Errorf("%w: missing artefact name", commonerrors.ErrUndefined)
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
		err = fmt.Errorf("%w: function to retrieve an artefact manager was not properly defined", commonerrors.ErrUndefined)
		return
	}

	err = filesystem.MkDir(outputDirectory)
	if err != nil {
		err = fmt.Errorf("%w: failed creating the output directory [%v] for job artefact: %v", commonerrors.ErrUnexpected, outputDirectory, err.Error())
		return
	}

	fileHasher, err := filesystem.NewFileHash(hashing.HashSha256)
	if err != nil {
		return
	}
	if any(artefactManager) == nil {
		err = fmt.Errorf("%w: missing artefact manager", commonerrors.ErrUndefined)
		return
	}

	artefactManagerName := artefactManager.GetName()
	if artefactManagerName == "" {
		err = fmt.Errorf("%w: missing artefact name", commonerrors.ErrUndefined)
		return
	}

	expectedSizePtr, ok := artefactManager.GetSizeOk()
	if !ok {
		err = fmt.Errorf("%w: could not fetch artefact's size from artefact's manager [%v]", commonerrors.ErrUndefined, artefactManagerName)
		return
	}
	expectedSize := *expectedSizePtr

	expectedHashPtr, ok := artefactManager.GetHashOk()
	if !ok {
		err = fmt.Errorf("%w: could not fetch artefact's hash from artefact's manager [%v]", commonerrors.ErrUndefined, artefactManagerName)
		return
	}
	expectedHash := *expectedHashPtr

	artefactFilename, artefactDestDir, err := determineArtefactDestination(outputDirectory, maintainTreeLocation, artefactManager)
	if err != nil {
		return
	}
	err = filesystem.MkDir(artefactDestDir)
	if err != nil {
		err = fmt.Errorf("%w: failed creating the output directory [%v] for job artefact: %v", commonerrors.ErrUnexpected, artefactDestDir, err.Error())
		return
	}

	artefact, resp, apierr := m.getArtefactContentFunc(ctx, jobName, artefactManagerName)
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
		if artefact != nil {
			_ = artefact.Close()
		}
	}()

	err = api.CheckAPICallSuccess(ctx, fmt.Sprintf("cannot fetch generated artefact [%v]", artefactFilename), resp, apierr)
	if err != nil {
		return
	}

	destination, err := filesystem.CreateFile(filepath.Join(artefactDestDir, artefactFilename))
	if err != nil {
		err = fmt.Errorf("%w: could not create a location to store generated artefact [%v]: %v", commonerrors.ErrUnexpected, artefactFilename, err.Error())
		return
	}
	defer func() { _ = destination.Close() }()

	actualSize, err := safeio.CopyDataWithContext(ctx, artefact, destination)
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
func (m *ArtefactManager[M, D, L, C]) DownloadJobArtefactFromLink(ctx context.Context, jobName string, outputDirectory string, artefactManagerItemLink D) error {
	return m.DownloadJobArtefactFromLinkWithTree(ctx, jobName, false, outputDirectory, artefactManagerItemLink)
}

func (m *ArtefactManager[M, D, L, C]) DownloadJobArtefactFromLinkWithTree(ctx context.Context, jobName string, maintainTreeLocation bool, outputDirectory string, artefactManagerItemLink D) (err error) {
	err = parallelisation.DetermineContextError(ctx)
	if err != nil {
		return
	}
	if m.getArtefactManagerFunc == nil || m.getArtefactContentFunc == nil {
		err = fmt.Errorf("%w: function to retrieve an artefact manager was not properly defined", commonerrors.ErrUndefined)
		return
	}
	if any(artefactManagerItemLink) == nil {
		err = fmt.Errorf("%w: missing artefact link", commonerrors.ErrUndefined)
		return
	}

	artefactManagerName := artefactManagerItemLink.GetName()
	artefactManager, resp, apierr := m.getArtefactManagerFunc(ctx, jobName, artefactManagerName)
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()
	err = api.CheckAPICallSuccess(ctx, fmt.Sprintf("cannot fetch artefact's manager [%+v]", artefactManager), resp, apierr)
	if err != nil {
		return
	}
	if resp != nil {
		_ = resp.Body.Close()
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
		err = fmt.Errorf("%w: function to retrieve artefact managers was not properly defined", commonerrors.ErrUndefined)
		return
	}
	clientPage, resp, apierr := m.getArtefactManagersFirstPageFunc(ctx, jobName)
	if resp != nil {
		_ = resp.Body.Close()
	}
	err = api.CheckAPICallSuccess(ctx, fmt.Sprintf("could not list artefact managers for job [%v]", jobName), resp, apierr)
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
		err = fmt.Errorf("%w: function to retrieve artefact managers was not properly defined", commonerrors.ErrUndefined)
		return
	}
	unwrappedPage := paginationUtils.UnwrapPage(currentPage)
	if unwrappedPage == nil {
		err = fmt.Errorf("%w: returned artefact managers page is empty", commonerrors.ErrUnexpected)
		return
	}
	page, ok := unwrappedPage.(C)
	if !ok {
		err = fmt.Errorf("%w: returned artefact managers page[%T] is not of the expected type [%v]", commonerrors.ErrUnexpected, currentPage, "*ArtefactManagerCollection")
		return
	}
	links, has := page.GetLinksOk()
	if !has {
		err = fmt.Errorf("%w: returned page of artefact managers has no links", commonerrors.ErrUnexpected)
		return
	}
	if !links.HasNext() {
		err = fmt.Errorf("%w: returned page of artefact managers has no `next` link", commonerrors.ErrUnexpected)
		return
	}
	link := links.GetNextP()
	clientPage, resp, apierr := m.getArtefactManagersFollowLinkFunc(ctx, link)
	if resp != nil {
		_ = resp.Body.Close()
	}
	err = api.CheckAPICallSuccess(ctx, fmt.Sprintf("could not follow `next` link [%v]", link), resp, apierr)
	if err == nil {
		nextPage = paginationUtils.ToPage(clientPage)
	}
	return
}

func (m *ArtefactManager[M, D, L, C]) DownloadAllJobArtefacts(ctx context.Context, jobName string, outputDirectory string) error {
	return m.DownloadAllJobArtefactsWithTree(ctx, jobName, false, outputDirectory)
}

func (m *ArtefactManager[M, D, L, C]) DownloadAllJobArtefactsWithTree(ctx context.Context, jobName string, maintainTreeStructure bool, outputDirectory string) (err error) {
	err = parallelisation.DetermineContextError(ctx)
	if err != nil {
		return
	}
	err = filesystem.MkDir(outputDirectory)
	if err != nil {
		err = fmt.Errorf("%w: failed creating the output directory [%v] for job artefacts: %v", commonerrors.ErrUnexpected, outputDirectory, err.Error())
		return
	}
	paginator, err := m.ListJobArtefacts(ctx, jobName)
	if err != nil {
		return
	}
	stop := paginator.Stop()
	defer stop()
	for {
		if !paginator.HasNext() {
			return
		}
		item, subErr := paginator.GetNext()
		if subErr != nil {
			err = fmt.Errorf("%w: failed getting information about job artefacts: %v", commonerrors.ErrUnexpected, subErr.Error())
			return
		}
		artefactLink, ok := item.(D)
		if ok {
			subErr = m.DownloadJobArtefactFromLinkWithTree(ctx, jobName, maintainTreeStructure, outputDirectory, artefactLink)
			if subErr != nil {
				err = subErr
				return
			}

		} else {
			artefactManager, ok := item.(M)
			if ok {
				subErr = m.DownloadJobArtefactWithTree(ctx, jobName, maintainTreeStructure, outputDirectory, artefactManager)
				if subErr != nil {
					err = subErr
					return
				}
			} else {
				err = fmt.Errorf("%w: the type of the response from service cannot be interpreted", commonerrors.ErrMarshalling)
				return
			}

		}

	}
}
