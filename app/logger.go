package app

import (
	"github.com/sisu-network/lib/log"
	tlog "github.com/tendermint/tendermint/libs/log"
)

var _ tlog.Logger = (*TendermintLogger)(nil)

// TendermintLogger is wrapper of logDNA and implement tendermint logger interface
type TendermintLogger struct {
	Inner *log.DNALogger
}

func NewTendermintLogger(dnaLogger *log.DNALogger) *TendermintLogger {
	return &TendermintLogger{Inner: dnaLogger}
}

func (l *TendermintLogger) Debug(msg string, keyvals ...interface{}) {
	// l.Inner.HighVerbose(msg, fmt.Sprintf("%v", keyvals))
}

func (l *TendermintLogger) Info(msg string, keyvals ...interface{}) {
	// l.Inner.HighVerbose(msg, fmt.Sprintf("%v", keyvals))
}

func (l *TendermintLogger) Error(msg string, keyvals ...interface{}) {
	// l.Inner.HighVerbose(msg, fmt.Sprintf("%v", keyvals))
}

func (l *TendermintLogger) With(keyvals ...interface{}) tlog.Logger {
	return l
}
