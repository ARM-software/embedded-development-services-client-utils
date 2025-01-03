/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// Package resource provides helpers related to API resources.
package resource

import "github.com/ARM-software/embedded-development-services-client/client"

// Mocks are generated using `go generate ./...`
// Add interfaces to the following command for a mock to be generated
//go:generate mockgen -destination=../mocks/mock_$GOPACKAGE.go -package=mocks github.com/ARM-software/embedded-development-services-client-utils/utils/$GOPACKAGE IResource

// IResource describe an API resource.
type IResource interface {
	client.IModel
}
