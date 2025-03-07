package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type TelegramSendBody struct {
	ChatID  int64  `json:"chat_id"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

type transportWithAuth struct {
	token string
}

func (t *transportWithAuth) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "ApiKey "+t.token)
	return http.DefaultTransport.RoundTrip(req)
}

type HttpClient struct {
	client          *http.Client
	parentAuthority string
}

func NewHttpClient(token string, parentAuthority string) *HttpClient {
	transport := transportWithAuth{token}
	client := &http.Client{Transport: &transport}

	return &HttpClient{
		client:          client,
		parentAuthority: parentAuthority,
	}
}

func (h *HttpClient) PollRequest(ctx context.Context) (Payload, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.parentAuthority+"/api/poll", nil)
	if err != nil {
		return Payload{}, err
	}

	resp, err := h.client.Do(req)
	if resp == nil || err != nil {
		return Payload{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Payload{}, err
	}

	payload := Payload{}
	err = json.Unmarshal(body, &payload)

	err = resp.Body.Close()
	if err != nil {
		return Payload{}, err
	}

	return payload, nil
}

func (h *HttpClient) TelegramSendRequest(
	ctx context.Context,
	nodeName string,
	chatID int64,
	message string,
) error {
	body := TelegramSendBody{
		ChatID:  chatID,
		Sender:  nodeName,
		Message: message,
	}

	out, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		h.parentAuthority+"/api/telegram-send",
		bytes.NewBuffer(out),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = h.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (h *HttpClient) TelegramSendToAllRequest(
	ctx context.Context,
	nodeName string,
	message string,
) error {
	body := TelegramSendBody{
		ChatID:  0,
		Sender:  nodeName,
		Message: message,
	}

	out, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		h.parentAuthority+"/api/telegram-send-to-all",
		bytes.NewBuffer(out),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = h.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
