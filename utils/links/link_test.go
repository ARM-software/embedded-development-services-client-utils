package links

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/links/linkstest"
)

func TestSerialiseLink(t *testing.T) {
	link, err := linkstest.NewFakeLink()
	require.NoError(t, err)
	assert.NotEmpty(t, SerialiseLink(link))

}
