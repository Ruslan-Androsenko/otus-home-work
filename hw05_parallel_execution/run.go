package hw05parallelexecution

import (
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var errorsCount, processedCount int
	tasksCount := len(tasks)

	// Создаем каналы для отправки задач на обработку и получения из них результатов
	doneCh := make(chan struct{})
	tasksCh := make(chan Task, n)
	errorsCh := make(chan error, tasksCount)

	// Создаем пулл исполнителей которые будут обрабатывать задачи
	for w := 0; w < n; w++ {
		go worker(doneCh, tasksCh, errorsCh)
	}

	// Запускаем обработку задач
	go executing(tasks, doneCh, tasksCh)

	// Проверяем количество полученных ошибок
	for err := range errorsCh {
		processedCount++

		if err != nil {
			errorsCount++
		}

		// Проверяем на лимит по ошибкам
		if m <= 0 && errorsCount > 0 || m > 0 && errorsCount >= m {
			close(doneCh)
			return ErrErrorsLimitExceeded
		}

		// Защищаемся от вечного цикла
		if processedCount >= tasksCount {
			close(doneCh)
			break
		}
	}
	close(errorsCh)

	return nil
}

// Обработка полученных задач.
func worker(doneCh <-chan struct{}, tasksCh <-chan Task, errorsCh chan<- error) {
	for {
		select {
		case <-doneCh:
			return
		case task := <-tasksCh:
			if task != nil {
				errorsCh <- task()
			}
		}
	}
}

// Добавление задач в обработку.
func executing(tasks []Task, doneCh <-chan struct{}, tasksCh chan<- Task) {
	defer close(tasksCh)

	for _, task := range tasks {
		select {
		case <-doneCh:
			return
		case tasksCh <- task:
		}
	}
}
