package app

import (
	"fmt"

	"github.com/sisu-network/lib/log"
	tlog "github.com/tendermint/tendermint/libs/log"
)

var _ tlog.Logger = (*Logger)(nil)

// Logger is wrapper of logDNA and implement tendermint logger interface
type Logger struct {
	dna *log.DNALogger
}

func NewTendermintLogger(dna *log.DNALogger) *Logger {
	return &Logger{dna: dna}
}

func (l *Logger) Debug(msg string, keyvals ...interface{}) {
	fmt.Println("come here debug")
	l.dna.Verbose(msg, fmt.Sprintf("%v", keyvals))
}

func (l *Logger) Info(msg string, keyvals ...interface{}) {
	fmt.Println("come here info")
	l.dna.Verbose(msg, fmt.Sprintf("%v", keyvals))
}

func (l *Logger) Error(msg string, keyvals ...interface{}) {
	fmt.Println("come here error")
	l.dna.Verbose(msg, fmt.Sprintf("%v", keyvals))
}

func (l *Logger) With(keyvals ...interface{}) tlog.Logger {
	fmt.Println("come here")
	return l
}
