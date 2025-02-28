package runner

import (
	"context"
	"errors"
	"fmt"
	"log"
	"miner-fetch/internal/handler/api"
	"net"
	"net/http"
	"time"
)

type HttpServer struct {
	CommonRunner
}

func NewHttpServer(runner CommonRunner) *HttpServer {
	ctxc, cancel := context.WithCancel(runner.ctx)
	runner.ctx = ctxc
	runner.cancel = cancel

	return &HttpServer{runner}
}

func (h *HttpServer) Start() {
	go func() {
		hs := api.NewHandler(h.cfg, h.s)

		mux := http.NewServeMux()

		mux.HandleFunc("GET /api/poll", hs.Poll)
		mux.HandleFunc("POST /api/telegram-send", hs.TelegramSend)

		handler := hs.AuthMiddleware(mux)

		server := &http.Server{
			Addr:    ":" + h.cfg.Port,
			Handler: handler,
			BaseContext: func(listener net.Listener) context.Context {
				return h.ctx
			},
		}

		go func() {
			if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}
		}()

		fmt.Printf("Server listening on port: %s\n", h.cfg.Port)

		<-h.ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			h.s.Logger.Log(err)
		}

		h.stopCh <- true
	}()
}

func (h *HttpServer) GetName() string {
	return "HttpServer"
}
