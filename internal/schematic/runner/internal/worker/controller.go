package worker

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/linkeunid/ligo"
	"{{.ModulePath}}/internal/usecase"
)

type Controller struct {
	workerUseCase *usecase.WorkerUseCase
	log           ligo.Logger
	cancel        context.CancelFunc
	running       atomic.Bool
}

func NewController(wuc *usecase.WorkerUseCase, log ligo.Logger) *Controller {
	return &Controller{
		workerUseCase: wuc,
		log:           log,
	}
}

func (c *Controller) Start() error {
	if c.running.Load() {
		return nil
	}

	c.log.Info("Worker controller starting worker")

	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	c.running.Store(true)

	go c.run(ctx)
	return nil
}

func (c *Controller) Stop() error {
	if !c.running.Load() {
		return nil
	}

	c.log.Info("Worker controller stopping worker")
	c.running.Store(false)
	if c.cancel != nil {
		c.cancel()
	}
	return nil
}

func (c *Controller) run(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	c.workerUseCase.Execute()

	for {
		select {
		case <-ctx.Done():
			c.log.Info("Worker controller: worker stopped")
			return
		case <-ticker.C:
			c.workerUseCase.Execute()
		}
	}
}
