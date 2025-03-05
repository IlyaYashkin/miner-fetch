package api

import (
	"encoding/json"
	"io"
	"miner-fetch/internal/config"
	"miner-fetch/internal/service"
	"net/http"
	"time"
)

type TelegramSendBody struct {
	ChatID  int64  `json:"chat_id"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

var pollTimeout = 30 * time.Second

type Handler struct {
	cfg config.Config
	s   *service.Service
}

func NewHandler(cfg config.Config, s *service.Service) *Handler {
	return &Handler{
		cfg: cfg,
		s:   s,
	}
}

func (h *Handler) Poll(w http.ResponseWriter, _ *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	clientChan := h.s.Polling.Subscribe()

	select {
	case payload := <-clientChan:
		msg, err := json.Marshal(payload)
		if err != nil {
			h.s.Logger.Log(err)
		}

		_, err = w.Write(msg)
		if err != nil {
			h.s.Logger.Log(err)
		}

	case <-time.After(pollTimeout):
		http.Error(w, "Timeout", http.StatusRequestTimeout)
	}

	h.s.Polling.Unsubscribe(clientChan)

	flusher.Flush()
}

func (h *Handler) TelegramSend(_ http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.s.Logger.Log(err)
	}

	resp := TelegramSendBody{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		h.s.Logger.Log(err)
	}

	err = h.s.TelegramSender.SendMessage(r.Context(), resp.ChatID, resp.Sender, resp.Message)
	if err != nil {
		h.s.Logger.Log(err)
	}
}
