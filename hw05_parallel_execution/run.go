package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, workersNumber, errorsNumber int) error {
	if errorsNumber <= 0 {
		return ErrErrorsLimitExceeded
	}

	var atomicErrorNumbers = int32(errorsNumber)
	var wg sync.WaitGroup

	taskChan := make(chan Task, len(tasks))
	doneChannel := make(chan bool, errorsNumber)
	var currentErrors int32

	go func() {
		defer func() {
			close(taskChan)
		}()
		for _, task := range tasks {
			taskChan <- task
		}
	}()

	for i := 0; i < workersNumber; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				if atomic.LoadInt32(&currentErrors) >= atomicErrorNumbers {
					return
				}
				err := task()
				if err != nil {
					atomic.AddInt32(&currentErrors, 1)
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(doneChannel)
	}()
	<-doneChannel

	if atomic.LoadInt32(&currentErrors) >= atomicErrorNumbers {
		return ErrErrorsLimitExceeded
	}
	return nil
}
