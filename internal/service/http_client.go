package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var PollTimeoutError = errors.New("poll timeout")
var NilResponseError = errors.New("nil response")

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

func (h *HttpClient) doRequest(req *http.Request) (*http.Response, error) {
	resp, err := h.client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, NilResponseError
	}

	if resp.StatusCode == http.StatusRequestTimeout {
		return resp, fmt.Errorf("request timed out")
	}

	if resp.StatusCode == http.StatusForbidden {
		return resp, fmt.Errorf("forbidden")
	}

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, err
}

func (h *HttpClient) PollRequest(ctx context.Context) (Payload, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.parentAuthority+"/api/poll", nil)
	if err != nil {
		return Payload{}, err
	}

	resp, err := h.doRequest(req)
	if !errors.Is(err, NilResponseError) && err != nil {
		return Payload{}, err
	}
	if !errors.Is(err, NilResponseError) && resp.StatusCode == http.StatusRequestTimeout {
		return Payload{}, PollTimeoutError
	}
	if err != nil {
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

	_, err = h.doRequest(req)
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

	_, err = h.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
