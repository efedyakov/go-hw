package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type Result struct {
	Message string
	Error   error
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
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
			reschan <- Result{"ошибка в горутине", err}
			return
		}
		reschan <- Result{"", nil}
	}

	index := 0
	for index < n {
		go dosome(tasks[index])
		index++
	}
	counterror := 0

	for index < len(tasks) {
		res := <-reschan
		if res.Error != nil {
			counterror++
		}
		if counterror > m {
			break
		}
		wg.Add(1)
		go dosome(tasks[index])
		index++
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
