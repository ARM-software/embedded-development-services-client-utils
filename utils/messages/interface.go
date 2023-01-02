/*
 * Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package messages

import (
	"context"
	"time"

	"github.com/ARM-software/golang-utils/utils/collection/pagination"
	"github.com/ARM-software/golang-utils/utils/logs"
)

// IMessage defines a generic service message.
type IMessage interface {
	// GetCtimeOk returns the creation time
	GetCtimeOk() (*time.Time, bool)
	// GetMessageOk returns the message string
	GetMessageOk() (*string, bool)
	// GetSeverityOk returns the message severity
	GetSeverityOk() (*string, bool)
	// GetSourceOk returns the message source.
	GetSourceOk() (*string, bool)
}

// IMessageLogger defines a logger dedicated to printing service messages.
type IMessageLogger interface {
	logs.Loggers
	LogMessage(msg IMessage)
	LogEmptyMessageError()
	LogMarshallingError(rawMessage *any)
	LogMessagesCollection(ctx context.Context, messagePaginator pagination.IGenericPaginator) error
}
