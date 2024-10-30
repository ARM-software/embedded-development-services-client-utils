/*
 * Copyright (C) 2020-2024 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package artefacts

import (
	"context"

	"github.com/ARM-software/embedded-development-services-client/client"
	"github.com/ARM-software/golang-utils/utils/collection/pagination"
)

//go:generate mockgen -destination=../mocks/mock_$GOPACKAGE.go -package=mocks github.com/ARM-software/embedded-development-services-client-utils/utils/$GOPACKAGE IArtefactManager

type IArtefactManager interface {
	// DownloadJobArtefactFromLink downloads a specific artefact into the output directory from a particular link.
	DownloadJobArtefactFromLink(ctx context.Context, jobName string, outputDirectory string, artefactManagerItemLink *client.HalLinkData) error
	// DownloadJobArtefact downloads a specific artefact into the output directory.
	DownloadJobArtefact(ctx context.Context, jobName string, outputDirectory string, artefactManager *client.ArtefactManagerItem) error
	// ListJobArtefacts lists all artefact managers associated with a particular job.
	ListJobArtefacts(ctx context.Context, jobName string) (pagination.IPaginatorAndPageFetcher, error)
	// DownloadAllJobArtefacts downloads all the artefacts produced for a particular job and puts them in an output directory.
	DownloadAllJobArtefacts(ctx context.Context, jobName string, outputDirectory string) error
}
