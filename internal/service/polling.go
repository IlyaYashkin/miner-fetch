package service

import (
	"sync"
)

type Payload struct {
	Command string
	ChatID  int64
}

type Polling struct {
	mu      sync.Mutex
	clients map[chan Payload]struct{}
}

func NewPolling() *Polling {
	return &Polling{
		clients: make(map[chan Payload]struct{}),
	}
}

func (p *Polling) Subscribe() chan Payload {
	p.mu.Lock()
	defer p.mu.Unlock()

	clientChan := make(chan Payload, 1)
	p.clients[clientChan] = struct{}{}

	return clientChan
}

func (p *Polling) Unsubscribe(client chan Payload) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.clients, client)
	close(client)
}

func (p *Polling) Send(payload Payload) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for client := range p.clients {
		client <- payload
	}
}
