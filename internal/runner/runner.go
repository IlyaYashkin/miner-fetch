package runner

import (
	"context"
	"miner-fetch/internal/config"
	"miner-fetch/internal/service"
)

type CommonRunner struct {
	ctx    context.Context
	cancel context.CancelFunc
	stopCh chan bool
	s      *service.Service
	cfg    config.Config
}

func NewCommonRunner(
	ctx context.Context,
	s *service.Service,
	cfg config.Config,
) CommonRunner {
	return CommonRunner{
		ctx:    ctx,
		stopCh: make(chan bool),
		s:      s,
		cfg:    cfg,
	}
}

func (c *CommonRunner) Stop() {
	c.cancel()

	<-c.stopCh
}
