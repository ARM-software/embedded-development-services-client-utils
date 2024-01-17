/*
 * Copyright (C) 2020-2024 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package logging

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-logr/logr"

	"github.com/ARM-software/embedded-development-services-client-utils/utils/resource"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/logs"
	"github.com/ARM-software/golang-utils/utils/reflection"
)

type ClientLogger struct {
	mLogger *logs.MultipleLoggerWithLoggerSource
}

func (c *ClientLogger) Close() error {
	return c.mLogger.Close()
}

func (c *ClientLogger) Check() error {
	return c.mLogger.Check()
}

func (c *ClientLogger) SetLogSource(source string) error {
	return c.mLogger.SetLogSource(source)
}

func (c *ClientLogger) SetLoggerSource(source string) error {
	return c.mLogger.SetLoggerSource(source)
}

func (c *ClientLogger) Log(output ...interface{}) {
	c.mLogger.Log(output...)
}

func (c *ClientLogger) LogError(err ...interface{}) {
	c.mLogger.LogError(err...)
}

func (c *ClientLogger) AppendLogger(l ...logr.Logger) error {
	return c.mLogger.AppendLogger(l...)
}

func (c *ClientLogger) Append(l ...logs.Loggers) error {
	return c.mLogger.Append(l...)
}

func (c *ClientLogger) LogRawError(err error) {
	c.LogErrorAndMessage(err, ": an error was encountered")
}

func (c *ClientLogger) LogErrorAndMessage(err error, format string, args ...interface{}) {
	var errorDescription string
	if format == "" {
		errorDescription = fmt.Sprint(args...)
	} else {
		errorDescription = fmt.Sprintf(format, args...)
	}
	if err == nil {
		err = errors.New(errorDescription)
	}
	c.LogError(err, ": ", errorDescription)
}

func (c *ClientLogger) LogErrorMessage(format string, args ...interface{}) {
	c.LogErrorAndMessage(nil, format, args...)
}

func (c *ClientLogger) LogInfo(format string, args ...interface{}) {
	var errorDescription string
	if format == "" {
		errorDescription = fmt.Sprint(args...)
	} else {
		errorDescription = fmt.Sprintf(format, args...)
	}
	c.Log(errorDescription)
}

func (c *ClientLogger) LogResource(r resource.IResource) {
	if r == nil {
		c.LogErrorAndMessage(commonerrors.ErrUndefined, "missing resource")
	} else {
		title, err := r.FetchTitle()
		if err != nil {
			c.LogErrorAndMessage(err, "could not retrieve resource's title")
			return
		}
		name, err := r.FetchName()
		if err != nil {
			c.LogErrorAndMessage(err, "could not retrieve resource's name")
			return
		}
		links, err := r.FetchLinks()
		if err != nil {
			c.LogErrorAndMessage(err, "could not retrieve resource's links [%v]", title)
			return
		}
		c.LogInfo("Resource (%v): %v [%v] ; affordances': %v", r.FetchType(), title, name, serialiseLink(links))
	}
}

func serialiseLink(links any) (str string) {
	if links == nil {
		str = "none"
		return
	}
	linkStr, err := json.Marshal(links)
	if err != nil {
		str = "unknown"
		return
	}
	str = string(linkStr)
	return
}

// NewClientLogger returns a logger for use in clients.
// if no default loggers are provided, the logger will be set to print to the standard output.
func NewClientLogger(loggerSource string, defaultLoggers ...logs.Loggers) (l ILogger, err error) {
	if loggerSource == "" {
		err = commonerrors.ErrNoLoggerSource
		return
	}
	multipleL, err := logs.NewMultipleLoggers(loggerSource, defaultLoggers...)
	if err != nil {
		return
	}
	multipleLogger, ok := multipleL.(*logs.MultipleLoggerWithLoggerSource)
	if !ok && multipleLogger == nil {
		err = fmt.Errorf("%w: multiple logger is not of the expected type", commonerrors.ErrUnexpected)
		return
	}
	l = &ClientLogger{
		mLogger: multipleLogger,
	}

	err = l.SetLoggerSource(loggerSource)
	return
}

// NewStandardClientLogger returns a typical client logger with logs written to a file if the logFilePath is set.
func NewStandardClientLogger(loggerSource string, logFilePath *string) (l ILogger, err error) {
	l, err = NewClientLogger(loggerSource)
	if err != nil || reflection.IsEmpty(logFilePath) {
		return
	}
	fileLogger, err := logs.NewFileLogger(*logFilePath, loggerSource)
	if err != nil {
		return
	}
	err = l.Append(fileLogger)
	return
}
