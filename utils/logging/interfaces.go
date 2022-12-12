package logging

import (
	"github.com/ARM-software/embedded-development-services-client-utils/utils/resource"
	"github.com/ARM-software/golang-utils/utils/logs"
)

type ILogger interface {
	logs.IMultipleLoggers
	LogRawError(err error)
	LogErrorAndMessage(err error, format string, args ...interface{})
	LogErrorMessage(format string, args ...interface{})
	LogInfo(format string, args ...interface{})
	LogResource(r resource.IResource)
}
