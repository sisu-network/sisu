package app

import (
	"fmt"

	"github.com/logdna/logdna-go/logger"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/config"
	tlog "github.com/tendermint/tendermint/libs/log"
)

var _ tlog.Logger = (*Logger)(nil)

// Logger is wrapper of logDNA and implement tendermint logger interface
type Logger struct {
	Inner *log.DNALogger
}

func NewTendermintLoggerIfHasSecret(cfg config.Config) *Logger {
	if len(cfg.LogDNA.Secret) == 0 {
		return nil
	}

	opts := logger.Options{
		App:           cfg.LogDNA.AppName,
		FlushInterval: cfg.LogDNA.FlushInterval.Duration,
		Hostname:      cfg.LogDNA.HostName,
		MaxBufferLen:  cfg.LogDNA.MaxBufferLen,
	}
	logDNA := log.NewDNALogger(cfg.LogDNA.Secret, opts)
	return &Logger{Inner: logDNA}
}

func (l *Logger) Debug(msg string, keyvals ...interface{}) {
	l.Inner.Verbose(msg, fmt.Sprintf("%v", keyvals))
}

func (l *Logger) Info(msg string, keyvals ...interface{}) {
	l.Inner.Verbose(msg, fmt.Sprintf("%v", keyvals))
}

func (l *Logger) Error(msg string, keyvals ...interface{}) {
	l.Inner.Verbose(msg, fmt.Sprintf("%v", keyvals))
}

func (l *Logger) With(keyvals ...interface{}) tlog.Logger {
	return l
}
