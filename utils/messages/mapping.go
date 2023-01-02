/*
 * Copyright (C) 2020-2022 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package messages

import (
	"github.com/ARM-software/embedded-development-services-client/client"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/reflection"
)

func convertRawMessageIntoIMessage(m *interface{}) (IMessage, error) {
	if reflection.IsEmpty(m) {
		return nil, commonerrors.ErrEmpty
	}
	raw := *m
	if messageItem, ok := raw.(client.NotificationMessageObject); ok {
		return &messageItem, nil
	}
	if messageItem, ok := raw.(client.MessageObject); ok {
		return &messageItem, nil
	}
	if messageItem, ok := raw.(IMessage); ok {
		return messageItem, nil
	}
	return nil, commonerrors.ErrMarshalling
}
