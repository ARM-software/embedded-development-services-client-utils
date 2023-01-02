/*
 * Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// Package resource provides helpers related to API resources.
package resource

import "github.com/ARM-software/embedded-development-services-client/client"

// IResource describe an API resource.
type IResource interface {
	client.IModel
}
