/*
 * Copyright (C) 2020-2024 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package linkstest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFakeLink(t *testing.T) {
	l, err := NewFakeLink()
	require.NoError(t, err)
	fmt.Println(l)
}
