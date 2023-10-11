package goroutines_utils

import (
	"fmt"
	"sync/atomic"
)

type WorkerPool struct {
	numWorkers uint8
	tasks      chan Runnable
	isAlive    atomic.Bool
}

func NewWorkerPool(numWorkers uint8) {
	pool := WorkerPool{
		numWorkers: numWorkers,
		tasks:      make(chan Runnable),
	}
}

func (w *WorkerPool) Add(runnables ...Runnable) error {
	if !w.isAlive.Load() {
		return fmt.Errorf("WorkerPool is dead")
	}
	for _, runnable := range runnables {
		w.tasks <- runnable
	}
	return nil
}
