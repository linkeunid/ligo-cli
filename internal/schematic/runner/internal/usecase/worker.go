package usecase

import (
	"github.com/linkeunid/ligo"
)

type WorkerUseCase struct {
	log ligo.Logger
}

func NewWorkerUseCase(log ligo.Logger) *WorkerUseCase {
	return &WorkerUseCase{log: log}
}

func (uc *WorkerUseCase) Execute() {
	uc.log.Info("Worker executing task")
	// Add your worker logic here
}
