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
)

// FormatMessage formats a job message.
func FormatMessage(msg IMessage) (s string, err error) {
	if msg == nil {
		err = fmt.Errorf("%w: missing message", commonerrors.ErrUndefined)
		return
	}
	var b strings.Builder
	if source, ok := msg.GetSourceOk(); ok {
		_, err = b.WriteString(fmt.Sprintf("[%s] ", *source))
		if err != nil {
			err = fmt.Errorf("%w: %v", commonerrors.ErrUnexpected, err.Error())
			return
		}
	}
	if ctime, ok := msg.GetCtimeOk(); ok {
		_, err = b.WriteString(fmt.Sprintf("(%s) ", *ctime))
		if err != nil {
			err = fmt.Errorf("%w: %v", commonerrors.ErrUnexpected, err.Error())
			return
		}
	}
	if severity, ok := msg.GetSeverityOk(); ok {
		_, err = b.WriteString(fmt.Sprintf("%s ", *severity))
		if err != nil {
			err = fmt.Errorf("%w: %v", commonerrors.ErrUnexpected, err.Error())
			return
		}
	}

	if message, ok := msg.GetMessageOk(); ok {
		if b.Len() > 0 {
			_, err = b.WriteString(fmt.Sprintf(": %s", *message))
			if err != nil {
				err = fmt.Errorf("%w: %v", commonerrors.ErrUnexpected, err.Error())
				return
			}
		} else {
			_, err = b.WriteString(*message)
			if err != nil {
				err = fmt.Errorf("%w: %v", commonerrors.ErrUnexpected, err.Error())
				return
			}
		}
	}
	s = b.String()
	return
}
