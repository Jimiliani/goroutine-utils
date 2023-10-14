package goroutines_utils

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type WorkerPool struct {
	numWorkers uint8
	tasks      chan Runnable
	wg         sync.WaitGroup
	isAlive    atomic.Bool
}

func NewWorkerPool(numWorkers uint8) *WorkerPool {
	pool := WorkerPool{
		numWorkers: numWorkers,
		tasks:      make(chan Runnable),
		wg:         sync.WaitGroup{},
	}
	var i uint8 = 0
	for ; i < numWorkers; i++ {
		pool.wg.Add(1)
		go func() {
			defer pool.wg.Done()
			for task := range pool.tasks {
				task.Run()
			}
		}()
	}

	pool.isAlive.Store(true)

	return &pool
}

func (w *WorkerPool) Add(runnables ...Runnable) (added uint64, err error) {
	if !w.isAlive.Load() {
		return 0, fmt.Errorf("WorkerPool is dead")
	}
	for _, runnable := range runnables {
		if runnable.acknowledge() {
			return added, fmt.Errorf("runnable %s already ran", runnable.GetID())
		}
		w.tasks <- runnable
		added += 1
	}
	return added, err
}

func (w *WorkerPool) Stop() {
	w.isAlive.Store(false)
	close(w.tasks)
	w.wg.Wait()
}
