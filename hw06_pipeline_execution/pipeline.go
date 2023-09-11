package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func controldone(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				out <- v
			}
		}
	}()
	return out
}

func dostage(in In, done In, stages ...Stage) Out {
	if len(stages) > 1 {
		return dostage(stages[0](controldone(in, done)), done, stages[1:]...)
	}
	return stages[0](controldone(in, done))
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Place your code here.
	inchan := make(Bi)
	go func() {
		defer close(inchan)
		for {
			select {
			case <-done:
				return

			case v, ok := <-in:
				if !ok {
					return
				}
				inchan <- v
			}
		}
	}()

	return dostage(inchan, done, stages...)
}
