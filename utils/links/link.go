/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// Package links provides utilities to deal with links
package links

import (
	"encoding/json"
)

// SerialiseLink serialises links into JSON strings
func SerialiseLink(links any) (str string) {
	if links == nil {
		return
	}
	linkStr, err := json.Marshal(links)
	if err != nil {
		return
	}
	str = string(linkStr)
	if str == "\"\"" {
		str = ""
	}
	return
}
