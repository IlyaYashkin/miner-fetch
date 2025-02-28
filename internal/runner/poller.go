package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"miner-fetch/internal/handler/api"
	"miner-fetch/internal/service"
	"net/http"
	"time"
)

type transportWithAuth struct {
	token string
}

func (t *transportWithAuth) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "ApiKey "+t.token)
	return http.DefaultTransport.RoundTrip(req)
}

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
		transport := transportWithAuth{p.cfg.AuthKey}
		client := &http.Client{Transport: &transport}

	L:
		for {
			select {
			case <-p.ctx.Done():
				p.stopCh <- true
				break L
			default:
				payload, err := p.pollRequest(p.ctx, client)

				if err != nil {
					p.s.Logger.Log(err)
				}

				if payload.Command != "" {
					message, err := p.s.Device.ExecuteCommand(payload.Command)
					if err != nil {
						p.s.Logger.Log(err)
					} else {
						err := p.telegramSendRequest(p.ctx, client, payload.ChatID, message)
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
) (service.Payload, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.cfg.ParentAuthority+"/api/poll", nil)
	if err != nil {
		return service.Payload{}, err
	}

	resp, err := client.Do(req)
	if resp != nil && resp.StatusCode == http.StatusRequestTimeout {
		return service.Payload{}, nil
	}

	if err != nil {
		return service.Payload{}, err
	}

	if resp == nil {
		return service.Payload{}, errors.New("nil response")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return service.Payload{}, err
	}

	payload := service.Payload{}
	err = json.Unmarshal(body, &payload)

	err = resp.Body.Close()
	if err != nil {
		return service.Payload{}, err
	}

	return payload, nil
}

func (p *Poller) telegramSendRequest(
	ctx context.Context,
	client *http.Client,
	chatID int64,
	message string,
) error {
	body := api.TelegramSendBody{
		ChatID:  chatID,
		Sender:  p.cfg.NodeName,
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
