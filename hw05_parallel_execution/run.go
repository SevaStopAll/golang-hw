package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func consume(taskChan <-chan Task, taskError chan error) {
	for {
		select {
		case task, ok := <-taskChan:
			if !ok {
				fmt.Println("Chan is closed")
				return
			}
			err := task()
			if err != nil {
				taskError <- err
			}
		case _, ok := <-taskError:
			if !ok {
				return
			}
		}
	}
}

func Run(tasks []Task, workersNumber, errorsNumber int) error {

	var wg sync.WaitGroup

	errorChannel := make(chan error, workersNumber)
	doneChannel := make(chan struct{}, errorsNumber)

	producerOwner := func() <-chan Task {
		tasksChan := make(chan Task, workersNumber)

		go func() {
			defer close(tasksChan)
			for _, task := range tasks {
				select {
				case <-doneChannel:
					return
				default:
					fmt.Println("Producer task", task)
					tasksChan <- task
				}
			}
		}()
		return tasksChan
	}

	producerChan := producerOwner()

	for i := 0; i < workersNumber; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("CONSUMER CLOSED")
			consume(producerChan, errorChannel)
		}()
	}

	go func() {
		wg.Wait()
		fmt.Println("wg waited")
		close(doneChannel)
	}()

	counter := 0
	for {
		select {
		case err, _ := <-errorChannel:
			if err != nil {
				counter++
				if counter >= errorsNumber {
					return ErrErrorsLimitExceeded
				}
			}
		case <-doneChannel:
			fmt.Println("Work is done")
			return nil
		}
	}
}
