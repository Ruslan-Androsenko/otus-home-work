package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type CheckError struct {
	DoneCh           chan struct{}
	TasksCh          chan Task
	wg               sync.WaitGroup
	ErrorsCounter    int32
	ProcessedCounter int32
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var (
		checker    CheckError
		maxErrors  = int32(m)
		tasksCount = int32(len(tasks))
	)

	// Создаем каналы для отправки задач на обработку, и сигнальный для завершения работы горутин
	checker.DoneCh = make(chan struct{})
	checker.TasksCh = make(chan Task, n)

	// Создаем пулл исполнителей которые будут обрабатывать задачи
	for w := 0; w < n; w++ {
		checker.wg.Add(1)
		go worker(&checker)
	}

	if maxErrors <= 0 {
		maxErrors = 1
	}

	// Запускаем обработку задач
	for _, task := range tasks {
		// Проверяем количество полученных ошибок
		if atomic.LoadInt32(&checker.ErrorsCounter) >= maxErrors {
			close(checker.DoneCh)
			close(checker.TasksCh)

			return ErrErrorsLimitExceeded
		}

		checker.TasksCh <- task
	}
	close(checker.TasksCh)

	for {
		// Проверяем что уже все задачи завершили свою работу
		if atomic.LoadInt32(&checker.ProcessedCounter) >= tasksCount {
			close(checker.DoneCh)
			break
		}
	}

	checker.wg.Wait()

	return nil
}

// Обработка полученных задач.
func worker(checker *CheckError) {
	defer checker.wg.Done()

	for {
		select {
		case <-checker.DoneCh:
			return
		case task, ok := <-checker.TasksCh:
			if ok {
				err := task()
				if err != nil {
					// Увеличиваем счетчик полученных ошибок
					atomic.AddInt32(&checker.ErrorsCounter, 1)
				}

				// Увеличиваем счетчик выполненных задач
				atomic.AddInt32(&checker.ProcessedCounter, 1)
			}
		}
	}
}
