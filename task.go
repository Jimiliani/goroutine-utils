package goroutines_utils

type Task[T any] struct {
	task       func() T
	IsFinished bool
	onFinish   chan struct{}
	// задача, которую можно уникально идентифицировать, определить жива ли она
	// получить результат, получить через канал событие о том, что таска завершила свою работу
}

func (t *Task[T]) Run() error {
	return nil
}
