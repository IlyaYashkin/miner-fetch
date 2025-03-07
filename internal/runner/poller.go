package runner

import (
	"context"
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
					if err != nil {
						p.s.Logger.Log(err)
					} else {
						err := p.s.HttpClient.TelegramSendRequest(p.ctx, p.cfg.NodeName, payload.ChatID, message)
						if err != nil {
							p.s.Logger.Log(err)
						}
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
