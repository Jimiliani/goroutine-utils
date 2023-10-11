package goroutines_utils

type TaskGroup struct {
	tasks []Runnable
}

func (t *TaskGroup) Run() error {
	return nil
}
