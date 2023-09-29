package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var (
		out         = in
		emptyOut    = make(Bi)
		countStages = len(stages)
	)
	close(emptyOut)

	// Если стейджи отсутствуют, сразу выходим из функции
	if countStages == 0 {
		return emptyOut
	}

	for _, stage := range stages {
		// Если done-канал не был передан, то не прослушиваем его
		if done != nil {
			for {
				select {
				case <-done:
					return emptyOut
				default:
					out = stage(out)
					break
				}
			}
		} else {
			out = stage(out)
		}
	}

	return out
}
