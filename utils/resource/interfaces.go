/*
 * Copyright (C) 2020-2022 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// Package resource provides helpers related to API resources.
package resource

// IResource describe an API resource.
type IResource interface {
	// GetLinks returns the resource links
	GetLinks() (any, error)
	GetName() (string, error)
	GetTitle() (string, error)
	GetType() string
}
