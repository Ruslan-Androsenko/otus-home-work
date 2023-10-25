package hw06pipelineexecution

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

func TestPipeline(t *testing.T) {
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	// Подготавливаем массив входных данных для большинства тестов.
	fiveElements := generateInputData(5)

	t.Run("simple case", func(t *testing.T) {
		in := sendInputData(fiveElements)

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, []string{"102", "104", "106", "108", "110"}, result)
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(fiveElements)-1)+int64(fault))
	})

	t.Run("done case", func(t *testing.T) {
		// Abort after 200ms
		abortDur := sleepPerStage * 2
		done := closeDoneChanelPerTime(abortDur)
		in := sendInputData(fiveElements)

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Len(t, result, 0)
		require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))
	})

	t.Run("empty stages", func(t *testing.T) {
		in := sendInputData(fiveElements)

		var emptyStages []Stage
		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, emptyStages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Len(t, result, 0)
		require.Less(t, int64(elapsed), int64(fault))
	})

	t.Run("additional stages", func(t *testing.T) {
		const countElements = 10

		data := generateInputData(countElements)
		in := sendInputData(data)

		countStages := len(stages)
		additionalStages := stages[:countStages-1]
		additionalStages = append(additionalStages,
			g("Multiplier (* 4)", func(v interface{}) interface{} { return v.(int) * 4 }),
			g("Adder (+ 512)", func(v interface{}) interface{} { return v.(int) + 512 }),
			g("Divider (/ 8)", func(v interface{}) interface{} { return v.(int) / 8 }),
			g("Subtraction (- 42)", func(v interface{}) interface{} { return v.(int) - 42 }),

			stages[countStages-1],
		)

		result := make([]string, 0, countElements)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, additionalStages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, []string{"73", "74", "75", "76", "77", "78", "79", "80", "81", "82"}, result)
		require.Less(t,
			int64(elapsed),
			// ~1.7s for processing 10 values in 8 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(additionalStages)+len(data)-1)+int64(fault))
	})

	t.Run("many data with done case", func(t *testing.T) {
		const countElements = 100

		// Abort after 2s
		abortDur := sleepPerStage * 20
		done := closeDoneChanelPerTime(abortDur)
		data := generateInputData(countElements)
		in := sendInputData(data)

		result := make([]string, 0, countElements)
		start := time.Now()
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		// Проверяем, что некоторые данные успели обработаться
		require.Less(t, 0, len(result))

		// Проверяем, что не все данные успели обработаться
		require.Less(t, len(result), countElements)
		require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))
	})
}

// Генерируем необходимое количество входных данных.
func generateInputData(counter int) []int {
	data := make([]int, counter)

	for i := 0; i < counter; i++ {
		data[i] = i + 1
	}

	return data
}

// Отправка входных данных во входной канал.
func sendInputData(data []int) In {
	in := make(Bi)

	go func() {
		for _, v := range data {
			in <- v
		}
		close(in)
	}()

	return in
}

// Закрыть сигнальный канал через необходимое время.
func closeDoneChanelPerTime(abortDur time.Duration) In {
	done := make(Bi)

	go func() {
		<-time.After(abortDur)
		close(done)
	}()

	return done
}
