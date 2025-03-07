package runner

import (
	"context"
	"errors"
	"miner-fetch/internal/service"
	"time"
)

var pollerSleep = 250 * time.Millisecond
var pollerSleepOnError = 10 * time.Second

type Poller struct {
	CommonRunner
}

func NewPoller(runner CommonRunner) *Poller {
	ctxc, cancel := context.WithCancel(runner.ctx)
	runner.ctx = ctxc
	runner.cancel = cancel

	return &Poller{runner}
}

func (p *Poller) Start() {
	go func() {
	L:
		for {
			select {
			case <-p.ctx.Done():
				p.stopCh <- true
				break L
			default:
				payload, err := p.s.HttpClient.PollRequest(p.ctx)

				if err != nil {
					p.s.Logger.Log(err)
					time.Sleep(pollerSleepOnError)
				}

				if payload.Command != "" {
					message, err := p.s.Device.ExecuteCommand(payload.Command)
					target := &service.CommandNotFoundError{}
					if err != nil && !errors.As(err, &target) {
						p.s.Logger.Log(err)
						continue
					} else if errors.As(err, &target) {
						message = err.Error()
					}

					err = p.s.HttpClient.TelegramSendRequest(p.ctx, p.cfg.NodeName, payload.ChatID, message)
					if err != nil {
						p.s.Logger.Log(err)
					}
				}

				time.Sleep(pollerSleep)
			}
		}
	}()
}

func (p *Poller) GetName() string {
	return "Poller"
}
