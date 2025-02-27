package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"miner-fetch/internal/handler/api"
	"net/http"
	"time"
)

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
		client := &http.Client{}

	L:
		for {
			select {
			case <-p.ctx.Done():
				p.stopCh <- true
				break L
			default:
				command, err := p.pollRequest(p.ctx, client)

				if err != nil {
					p.s.Logger.Log(err)
				}

				if command != "" {
					message, err := p.s.Device.ExecuteCommand(command)
					if err != nil {
						p.s.Logger.Log(err)
					} else {
						fmt.Println(message)
						err := p.telegramSendRequest(p.ctx, client, message)
						if err != nil {
							p.s.Logger.Log(err)
						}
					}
				}

				time.Sleep(250 * time.Millisecond)
			}
		}
	}()
}

func (p *Poller) GetName() string {
	return "Poller"
}

func (p *Poller) pollRequest(
	ctx context.Context,
	client *http.Client,
) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.cfg.ParentAuthority+"/api/poll", nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if resp != nil && resp.StatusCode == http.StatusRequestTimeout {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	if resp == nil {
		return "", errors.New("nil response")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = resp.Body.Close()
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (p *Poller) telegramSendRequest(
	ctx context.Context,
	client *http.Client,
	message string,
) error {
	body := api.TelegramSendBody{
		ChatID:  123123,
		Sender:  "LOPATINA",
		Message: message,
	}

	out, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.cfg.ParentAuthority+"/api/telegram-send", bytes.NewBuffer(out))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
