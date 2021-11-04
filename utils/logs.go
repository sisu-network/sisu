package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
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
var loggerLock = &sync.RWMutex{}

func SetLogger(logger Logger) {
	loggerLock.Lock()
	defer loggerLock.Unlock()

	globalLogger = logger
}

func SetLogLevel(level int) {
	loggerLock.Lock()
	defer loggerLock.Unlock()

	logLevel = level
}

func getLogger() Logger {
	loggerLock.RLock()

	if globalLogger == nil {
		// Release the lock so that it could be locked again in the set functions
		loggerLock.RUnlock()

		SetLogger(newDefaultLogger())
		SetLogLevel(LOG_LEVEL_INFO)
	} else {
		loggerLock.RUnlock()
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
	errorColor   *color.Color
	verboseColor *color.Color
	warningColor *color.Color
	debugColor   *color.Color
}

func newDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		errorColor:   color.New(color.FgRed),
		verboseColor: color.New(color.FgCyan),
		warningColor: color.New(color.FgYellow),
		debugColor:   color.New(color.FgMagenta),
	}
}

func (logger *DefaultLogger) LogInfo(a ...interface{}) {
	logger.printWithTime(nil, a...)
}

func (logger *DefaultLogger) LogDebug(a ...interface{}) {
	logger.printWithTime(logger.debugColor, a...)
}

func (logger *DefaultLogger) LogWarn(a ...interface{}) {
	logger.printWithTime(logger.warningColor, a...)
}

func (logger *DefaultLogger) LogError(a ...interface{}) {
	logger.printWithTime(logger.errorColor, a...)
}

func (logger *DefaultLogger) LogVerbose(a ...interface{}) {
	logger.printWithTime(logger.verboseColor, a...)
}

func (logger *DefaultLogger) LogCritical(a ...interface{}) {
	logger.printWithTime(nil, a...)
}

func (logger *DefaultLogger) printWithTime(c *color.Color, a ...interface{}) {
	now := time.Now().Format("15:04:05.00")

	var m []interface{}
	m = append(m, now)
	m = append(m, a...)

	if c == nil {
		fmt.Println(m...)
	} else {
		c.Println(m...)
	}
}
