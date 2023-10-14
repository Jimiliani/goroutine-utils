package goroutines_utils

import (
	"sync"
	"sync/atomic"
)

type Task[T any] struct {
	id            string // не уникально
	task          func() (T, error)
	isFinished    atomic.Bool
	isRan         atomic.Bool
	mu            sync.RWMutex
	result        T
	error         error
	onFinishChans []chan struct{}
	// задача, которую можно уникально идентифицировать, определить жива ли она
	// получить результат, получить через канал событие о том, что таска завершила свою работу
}

func NewTask[T any](id string, task func() (T, error)) Task[T] {
	return Task[T]{
		id:         id,
		task:       task,
		isFinished: atomic.Bool{},
		isRan:      atomic.Bool{},
	}
}

func (t *Task[T]) Run() {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			t.mu.Lock()
			t.error = err
			t.mu.Unlock()
		}
	}()

	defer func() {
		t.isFinished.Store(true)
		go func() {
			for _, finishChan := range t.onFinishChans {
				finishChan <- struct{}{}
				close(finishChan)
			}
		}()
	}()

	res, err := t.task()
	t.mu.Lock()
	t.error = err
	t.result = res
	t.mu.Unlock()
}

func (t *Task[T]) acknowledge() bool {
	return t.isRan.Swap(true)
}

func (t *Task[T]) GetID() string {
	return t.id
}

func (t *Task[T]) GetResult() T {
	t.mu.RLock()
	res := t.result
	t.mu.RUnlock()
	return res
}

func (t *Task[T]) GetError() error {
	t.mu.RLock()
	err := t.error
	t.mu.RUnlock()
	return err
}

func (t *Task[T]) IsFinished() bool {
	return t.isFinished.Load()
}

func (t *Task[T]) OnFinish() <-chan struct{} {
	c := make(chan struct{}, 1)
	if t.isFinished.Load() {
		c <- struct{}{}
		close(c)
	} else {
		t.onFinishChans = append(t.onFinishChans, c)
		if t.isFinished.Load() {
			c <- struct{}{}
			close(c)
		}
	}
	return c
}
