/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// Package messages provides utilities to log messages
package messages

import (
	"fmt"
	"strings"

	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/field"
)

type FormatterOptions struct {
	source    bool
	timestamp bool
	severity  bool
}

type FormatterOption func(*FormatterOptions) *FormatterOptions

// WithSource includes information about the source of the message.
var WithSource FormatterOption = func(o *FormatterOptions) *FormatterOptions {
	if o == nil {
		return o
	}
	o.source = true
	return o
}

// WithTimeStamp includes the timestamp of the message.
var WithTimeStamp FormatterOption = func(o *FormatterOptions) *FormatterOptions {
	if o == nil {
		return o
	}
	o.timestamp = true
	return o
}

// WithSeverity includes the severity of the message.
var WithSeverity FormatterOption = func(o *FormatterOptions) *FormatterOptions {
	if o == nil {
		return o
	}
	o.severity = true
	return o
}

type MessageFormatter struct {
	options FormatterOptions
}

func (f *MessageFormatter) FormatMessage(msg IMessage) (s string, err error) {
	if msg == nil {
		err = commonerrors.UndefinedVariable("message")
		return
	}
	var b strings.Builder
	if f.options.source {
		if source, ok := msg.GetSourceOk(); ok {
			_, err = b.WriteString(fmt.Sprintf("[%s] ", field.OptionalString(source, "")))
			if err != nil {
				err = commonerrors.WrapError(commonerrors.ErrUnexpected, err, "")
				return
			}
		}
	}
	if f.options.timestamp {
		if ctime, ok := msg.GetCtimeOk(); ok {
			_, err = b.WriteString(fmt.Sprintf("(%s) ", *ctime))
			if err != nil {
				err = commonerrors.WrapError(commonerrors.ErrUnexpected, err, "")
				return
			}
		}
	}
	if f.options.severity {
		if severity, ok := msg.GetSeverityOk(); ok {
			_, err = b.WriteString(fmt.Sprintf("%s ", field.OptionalString(severity, "")))
			if err != nil {
				err = commonerrors.WrapError(commonerrors.ErrUnexpected, err, "")
				return
			}
		}
	}

	if message, ok := msg.GetMessageOk(); ok {
		if b.Len() > 0 {
			_, err = b.WriteString(fmt.Sprintf(": %s", field.OptionalString(message, "")))
			if err != nil {
				err = commonerrors.WrapError(commonerrors.ErrUnexpected, err, "")
				return
			}
		} else {
			_, err = b.WriteString(field.OptionalString(message, ""))
			if err != nil {
				err = commonerrors.WrapError(commonerrors.ErrUnexpected, err, "")
				return
			}
		}
	}
	s = b.String()
	return
}

// NewMessageFormatter creates a formatter for messages
func NewMessageFormatter(option ...FormatterOption) *MessageFormatter {
	options := &FormatterOptions{}
	for i := range option {
		options = option[i](options)
	}
	return &MessageFormatter{options: *options}
}

// DefaultMessageFormatter returns a default message formatter.
func DefaultMessageFormatter() *MessageFormatter {
	return NewMessageFormatter(WithSource, WithTimeStamp, WithSeverity)
}

// FormatMessageWithOptions formats a job message with formatting options.
func FormatMessageWithOptions(msg IMessage, option ...FormatterOption) (s string, err error) {
	return NewMessageFormatter(option...).FormatMessage(msg)
}

// FormatMessage formats a job message.
func FormatMessage(msg IMessage) (s string, err error) {
	return FormatMessageWithOptions(msg, WithSource, WithTimeStamp, WithSeverity)
}
