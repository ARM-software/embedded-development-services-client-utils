/*
 * Copyright (C) 2020-2024 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package messages

import (
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/reflection"
)

func convertRawMessageIntoIMessage(m interface{}) (IMessage, error) {
	if reflection.IsEmpty(m) {
		return nil, commonerrors.ErrEmpty
	}
	if messageItem, ok := m.(IMessage); ok {
		return messageItem, nil
	}
	return nil, commonerrors.ErrMarshalling
}
