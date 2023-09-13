package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, atomic.LoadInt32(&runTasksCount), int32(workersCount+maxErrorsCount),
			"extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, atomic.LoadInt32(&runTasksCount), int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("stable tasks with one random error and with negative limit errors", func(t *testing.T) {
		tasksCount := 70
		tasks := make([]Task, 0, tasksCount)
		rand.Seed(time.Now().UnixNano())
		min := 10
		randomTaskWithErr := rand.Intn((tasksCount - min + 1) + min)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(120))
			sumTime += taskSleep

			if i == randomTaskWithErr {
				tasks = append(tasks, func() error {
					time.Sleep(taskSleep)
					atomic.AddInt32(&runTasksCount, 1)

					return fmt.Errorf("error from task %d", i)
				})
			} else {
				tasks = append(tasks, func() error {
					time.Sleep(taskSleep)
					atomic.AddInt32(&runTasksCount, 1)

					return nil
				})
			}
		}

		workersCount := 5
		maxErrorsCount := -1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)

		require.Less(t, atomic.LoadInt32(&runTasksCount), int32(tasksCount),
			"not all tasks were accepted due to an accidental error")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("tasks that have a maximum of zero errors", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(170)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 4
		maxErrorsCount := 0
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, atomic.LoadInt32(&runTasksCount), int32(workersCount+maxErrorsCount),
			"extra tasks were started")
	})

	t.Run("tasks in which mistakes are made", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(200)))
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 4
		maxErrorsCount := 0
		err := Run(tasks, workersCount, maxErrorsCount)

		require.NoError(t, err)
		require.Equal(t, atomic.LoadInt32(&runTasksCount), int32(tasksCount), "not all tasks were completed")
	})
}
