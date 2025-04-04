/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package artefacts

import (
	"context"

	"github.com/ARM-software/golang-utils/utils/collection/pagination"
)

//go:generate mockgen -destination=../mocks/mock_$GOPACKAGE.go -package=mocks github.com/ARM-software/embedded-development-services-client-utils/utils/$GOPACKAGE IArtefactManager

type IArtefactManager[
	M Manager,
	D LinkData,
] interface {
	// DownloadJobArtefactFromLink downloads a specific artefact into the output directory from a particular link. The artefact will be placed at the root of the output directory.
	DownloadJobArtefactFromLink(ctx context.Context, jobName string, outputDirectory string, artefactManagerItemLink D) error
	// DownloadJobArtefactFromLinkWithTree downloads a specific artefact into the output directory from a particular link.
	// maintainTreeLocation specifies whether the artefact will be placed in a tree structure or if it will be flat.
	DownloadJobArtefactFromLinkWithTree(ctx context.Context, jobName string, maintainTreeLocation bool, outputDirectory string, artefactManagerItemLink D) error
	// DownloadJobArtefact downloads a specific artefact into the output directory. The artefact will be placed at the root of the output directory.
	DownloadJobArtefact(ctx context.Context, jobName string, outputDirectory string, artefactManager M) error
	// DownloadJobArtefactWithTree downloads a specific artefact into the output directory.
	// maintainTreeLocation specifies whether the artefact will be placed in a tree structure or if it will be flat.
	DownloadJobArtefactWithTree(ctx context.Context, jobName string, maintainTreeLocation bool, outputDirectory string, artefactManager M) error
	// ListJobArtefacts lists all artefact managers associated with a particular job.
	ListJobArtefacts(ctx context.Context, jobName string) (pagination.IPaginatorAndPageFetcher, error)
	// DownloadAllJobArtefacts downloads all the artefacts produced for a particular job and puts them in an output directory as a flat list.
	DownloadAllJobArtefacts(ctx context.Context, jobName string, outputDirectory string) error
	// DownloadAllJobArtefactsWithTree downloads all the artefacts produced for a particular job and puts them in an output directory.
	// maintainTreeStructure specifies whether to keep the tree structure of the artefacts or not in the output directory.
	DownloadAllJobArtefactsWithTree(ctx context.Context, jobName string, maintainTreeStructure bool, outputDirectory string) error
}

type LinkData interface {
	comparable
	GetName() string
	GetHref() string
}

type Links[D LinkData] interface {
	comparable
	GetNextP() D
	HasNext() bool
}
type Collection[D LinkData, L Links[D]] interface {
	comparable
	pagination.IStaticPage
	GetLinksOk() (L, bool)
}

type Manager interface {
	comparable
	GetName() string
	GetTitle() string
	HasTitle() bool
	GetHashOk() (*string, bool)
	GetSizeOk() (*int64, bool)
	GetExtraMetadata() map[string]string
	HasExtraMetadata() bool
}
