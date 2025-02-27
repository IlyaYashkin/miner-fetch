package runner

import (
	"context"
)

type Logger struct {
	CommonRunner
}

func NewLogger(runner CommonRunner) *Logger {
	ctxc, cancel := context.WithCancel(runner.ctx)
	runner.ctx = ctxc
	runner.cancel = cancel

	return &Logger{runner}
}

func (l *Logger) Start() {
	go func() {
		go l.s.Logger.Start()

		<-l.ctx.Done()

		l.s.Logger.Stop()

		l.stopCh <- true
	}()
}

func (l *Logger) GetName() string {
	return "Logger"
}
