package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded     = errors.New("errors limit exceeded")
	ErrErrorsNotExistGoRoutine = errors.New("not exist goroutine")
)

type Task func() error

type Result struct {
	// Message string
	Error error
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return ErrErrorsNotExistGoRoutine
	}
	if n > len(tasks) {
		n = len(tasks)
	}
	reschan := make(chan Result, n)

	wg := sync.WaitGroup{}
	wg.Add(n)

	dosome := func(task Task) {
		defer wg.Done()
		err := task()
		if err != nil {
			reschan <- Result{err}
			return
		}
		reschan <- Result{nil}
	}

	for i := 0; i < n; i++ {
		go dosome(tasks[i])
	}
	counterror := 0

	for i := n; i < len(tasks); i++ {
		res := <-reschan
		if res.Error != nil {
			counterror++
		}
		if counterror > m {
			break
		}
		wg.Add(1)
		go dosome(tasks[i])
	}

	wg.Wait()
	// получим результаты последних тасков
	close(reschan)
	for i := 0; i < n; i++ {
		res, ok := <-reschan
		if !ok {
			break
		}
		if res.Error != nil {
			counterror++
		}
	}

	if counterror > m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
