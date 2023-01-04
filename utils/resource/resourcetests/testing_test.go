/*
 * Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package resourcetests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewResourceTest(t *testing.T) {
	l, err := NewMockResource()
	require.NoError(t, err)
	fmt.Println(l)
}
