package service

import (
	"sync"
)

type Polling struct {
	mu      sync.Mutex
	clients map[chan string]struct{}
}

func NewPolling() *Polling {
	return &Polling{
		clients: make(map[chan string]struct{}),
	}
}

func (p *Polling) Subscribe() chan string {
	p.mu.Lock()
	defer p.mu.Unlock()

	clientChan := make(chan string, 1)
	p.clients[clientChan] = struct{}{}

	return clientChan
}

func (p *Polling) Unsubscribe(client chan string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.clients, client)
	close(client)
}

func (p *Polling) Send(msg string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for client := range p.clients {
		client <- msg
	}
}
