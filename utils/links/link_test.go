/*
 * Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package links

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/links/linkstest"
)

func TestSerialiseLink(t *testing.T) {
	assert.Empty(t, SerialiseLink(nil))
	assert.Empty(t, SerialiseLink(""))
	link, err := linkstest.NewFakeLink()
	require.NoError(t, err)
	assert.NotEmpty(t, SerialiseLink(link))
	assert.NotEmpty(t, SerialiseLink(&link))

}
