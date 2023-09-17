package hw05parallelexecution

import (
	"errors"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type CheckError struct {
	DoneCh       chan struct{}
	ErrorsCh     chan error
	TasksCh      chan Task
	HasErrClosed uint32
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var (
		errorsCount, processedCount int
		tasksCount                  = len(tasks)
		checker                     CheckError
	)

	// Создаем каналы для отправки задач на обработку и получения из них результатов
	checker.DoneCh = make(chan struct{})
	checker.ErrorsCh = make(chan error, n)
	checker.TasksCh = make(chan Task, n)

	// Закрываем канал с ошибками при выходе из основной функции
	defer close(checker.ErrorsCh)

	// Создаем пулл исполнителей которые будут обрабатывать задачи
	for w := 0; w < n; w++ {
		go worker(&checker)
	}

	// Запускаем обработку задач
	go executing(tasks, &checker)

	if m <= 0 {
		m = 1
	}

	// Проверяем количество полученных ошибок
	for err := range checker.ErrorsCh {
		processedCount++

		if err != nil {
			errorsCount++
		}

		// Проверяем на лимит по ошибкам
		if errorsCount >= m {
			close(checker.DoneCh)
			atomic.AddUint32(&checker.HasErrClosed, 1)

			return ErrErrorsLimitExceeded
		}

		// Защищаемся от вечного цикла
		if processedCount >= tasksCount {
			close(checker.DoneCh)
			atomic.AddUint32(&checker.HasErrClosed, 1)

			break
		}
	}

	return nil
}

// Обработка полученных задач.
func worker(checker *CheckError) {
	for {
		select {
		case <-checker.DoneCh:
			return
		case task, ok := <-checker.TasksCh:
			if ok {
				err := task()

				// Если канал с ошибками еще не закрыт, то добавляем в него результат выполнения задачи
				if atomic.LoadUint32(&checker.HasErrClosed) == 0 {
					checker.ErrorsCh <- err
				}
			}
		}
	}
}

// Добавление задач в обработку.
func executing(tasks []Task, checker *CheckError) {
	defer close(checker.TasksCh)

	for _, task := range tasks {
		select {
		case <-checker.DoneCh:
			return
		case checker.TasksCh <- task:
		}
	}
}
