package utils

import (
	"fmt"
	"time"
)

type Logger interface {
	LogCritical(a ...interface{})
	LogError(a ...interface{})
	LogWarn(a ...interface{})
	LogInfo(a ...interface{})
	LogVerbose(a ...interface{})
	LogDebug(a ...interface{})
}

const (
	LOG_LEVEL_DEBUG    = 1
	LOG_LEVEL_VERBOSE  = 10
	LOG_LEVEL_INFO     = 20
	LOG_LEVEL_WARN     = 30
	LOG_LEVEL_ERROR    = 40
	LOG_LEVEL_CRITICAL = 50
)

var globalLogger Logger
var logLevel int

func SetLogger(logger Logger) {
	globalLogger = logger
}

func SetLogLevel(level int) {
	logLevel = level
}

func getLogger() Logger {
	if globalLogger == nil {
		SetLogger(&DefaultLogger{})
		SetLogLevel(LOG_LEVEL_INFO)
	}

	return globalLogger
}

func LogInfo(a ...interface{}) {
	logger := getLogger()
	logger.LogInfo(a...)
}

func LogDebug(a ...interface{}) {
	logger := getLogger()
	logger.LogDebug(a...)
}

func LogWarn(a ...interface{}) {
	logger := getLogger()
	logger.LogWarn(a...)
}

func LogError(a ...interface{}) {
	logger := getLogger()
	logger.LogError(a...)
}

func LogVerbose(a ...interface{}) {
	logger := getLogger()
	logger.LogVerbose(a...)
}

func LogCritical(a ...interface{}) {
	logger := getLogger()
	logger.LogCritical(a...)
}

// Default logging implementation. You can replace this logging module by another implementation
// of Logger interface.

type DefaultLogger struct {
}

func (logger DefaultLogger) LogInfo(a ...interface{}) {
	logger.printWithTime(a...)
}

func (logger DefaultLogger) LogDebug(a ...interface{}) {
	logger.printWithTime(a...)
}

func (logger DefaultLogger) LogWarn(a ...interface{}) {
	logger.printWithTime(a...)
}

func (logger DefaultLogger) LogError(a ...interface{}) {
	logger.printWithTime(a...)
}

func (logger DefaultLogger) LogVerbose(a ...interface{}) {
	logger.printWithTime(a...)
}

func (logger DefaultLogger) LogCritical(a ...interface{}) {
	logger.printWithTime(a...)
}

func (logger DefaultLogger) printWithTime(a ...interface{}) {
	now := time.Now().Format("15:04:05.00")

	var m []interface{}
	m = append(m, now)
	m = append(m, a...)

	fmt.Println(m...)
}
