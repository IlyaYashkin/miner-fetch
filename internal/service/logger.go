package service

import (
	"log"
)

type Logger struct {
	ch        chan error
	isStarted bool
}

func NewLogger() *Logger {
	return &Logger{
		ch:        make(chan error, 100),
		isStarted: false,
	}
}

func (l *Logger) Log(msg error) {
	if l.isStarted {
		l.ch <- msg
	}
}

func (l *Logger) Start() {
	l.isStarted = true

	for msg := range l.ch {
		log.Println(msg)
	}
}

func (l *Logger) Stop() {
	l.isStarted = false

	close(l.ch)
}
