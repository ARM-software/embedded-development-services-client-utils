package artefacts

import (
	"context"

	"github.com/ARM-software/embedded-development-services-client/client"
)

//go:generate mockgen -destination=../mocks/mock_$GOPACKAGE.go -package=mocks github.com/ARM-software/embedded-development-services-client-utils/utils/$GOPACKAGE IArtefactManager

type IArtefactManager interface {
	DownloadJobArtefact(ctx context.Context, jobName string, outputDirectory string, artefactManagerItem client.HalLinkData) (err error)
}
