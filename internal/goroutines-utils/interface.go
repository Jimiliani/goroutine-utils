package goroutines_utils

type Runnable interface {
	Run()
	GetID() string
	acknowledge() bool
}
