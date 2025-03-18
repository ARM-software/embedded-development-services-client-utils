/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package client

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"

	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/commonerrors/errortest"
	"github.com/ARM-software/golang-utils/utils/http"
	"github.com/ARM-software/golang-utils/utils/logs/logrimp"
)

func TestNewClient(t *testing.T) {
	cfg := http.DefaultHTTPRequestWithAuthorisationConfigurationEnforced(faker.Sentence())
	c, err := NewClient(cfg, logr.Discard(), nil)
	errortest.RequireError(t, err, commonerrors.ErrNoLogger)
	assert.Empty(t, c)
	c, err = NewClient(cfg, logrimp.NewStdOutLogr(), nil)
	errortest.RequireError(t, err, commonerrors.ErrInvalid)
	assert.Empty(t, c)
	c, err = NewClient(nil, logrimp.NewStdOutLogr(), nil)
	errortest.RequireError(t, err, commonerrors.ErrUndefined)
	assert.Empty(t, c)
}
