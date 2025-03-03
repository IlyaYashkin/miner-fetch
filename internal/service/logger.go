package service

import (
	"fmt"
	"log"
	"runtime/debug"
)

type Logger struct {
	debug bool
}

func NewLogger(debug bool) *Logger {
	return &Logger{
		debug: debug,
	}
}

func (l *Logger) Log(msg error) {
	if l.debug {
		msg = fmt.Errorf("%w\nStack trace:\n%s", msg, debug.Stack())
	}

	log.Println(msg)
}
