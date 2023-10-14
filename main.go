package main

import (
	"fmt"
	"time"

	gutils "goroutines-utils/internal/goroutines-utils"
)

func main() {
	pool := gutils.NewWorkerPool(5)
	defer pool.Stop()

	tasks := make([]gutils.Runnable, 0, 10)
	var t1 *gutils.Task[int]
	for i := 0; i < 10; i++ {
		t := gutils.NewTask[int](
			fmt.Sprintf("t%d", i),
			func() (int, error) {
				// simulate http request
				time.Sleep(1 * time.Second)
				return 111, nil
			},
		)
		t1 = &t
		tasks = append(
			tasks,
			&t,
		)
	}

	_, err := pool.Add(tasks...)
	if err != nil {
		fmt.Println(fmt.Errorf("unable to add task to pool: %w", err))
		return
	}

	t := time.NewTicker(time.Millisecond * 500)
	defer t.Stop()

	for {
		select {
		case <-t1.OnFinish():
			fmt.Println("finish!")
			fmt.Println(t1.GetResult())
			fmt.Println(t1.GetError())
			return
		case <-t.C:
			fmt.Println("tick")
		}
	}

}
