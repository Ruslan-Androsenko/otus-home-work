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

	t.Run("simple case", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

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
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})

	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

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
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

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
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		countStages := len(stages)
		additionalStages := stages[:countStages-1]
		additionalStages = append(additionalStages,
			g("Multiplier (* 4)", func(v interface{}) interface{} { return v.(int) * 4 }),
			g("Adder (+ 512)", func(v interface{}) interface{} { return v.(int) + 512 }),
			g("Divider (/ 8)", func(v interface{}) interface{} { return v.(int) / 8 }),
			g("Subtraction (- 42)", func(v interface{}) interface{} { return v.(int) - 42 }),

			stages[countStages-1],
		)

		result := make([]string, 0, 10)
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
}
