/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// Package messages provides utilities to handle messages
package messages

import (
	"context"
	"time"

	"github.com/ARM-software/golang-utils/utils/collection/pagination"
	"github.com/ARM-software/golang-utils/utils/logs"
)

// Mocks are generated using `go generate ./...`
// Add interfaces to the following command for a mock to be generated
//go:generate go tool mockgen -destination=../mocks/mock_$GOPACKAGE.go -package=mocks github.com/ARM-software/embedded-development-services-client-utils/utils/$GOPACKAGE IMessage,IMessageLogger

const (
	// DefaultMessagesPrintingFrequency describes the default frequency at which messages are printed
	DefaultMessagesPrintingFrequency = 100 * time.Millisecond

	// DefaultMessageFetchingBackoff describes the fixed backoff duration before more messages are fetched.
	DefaultMessageFetchingBackoff = 10 * time.Millisecond

	// DefaultStreamExhaustionGracePeriod describes the grace period which should happen when expecting message stream exhaustion.
	DefaultStreamExhaustionGracePeriod = time.Second
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
