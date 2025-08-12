package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func drainStage(stageOut In) {
	for v := range stageOut {
		_ = v // for lint
	}
}

func processStage(done In, stageOut In, out Bi) {
	defer close(out)
	for {
		select {
		case <-done:
			go drainStage(stageOut)
			return
		case value, ok := <-stageOut:
			if !ok {
				return
			}
			select {
			case out <- value:
			case <-done:
				go drainStage(stageOut)
				return
			}
		}
	}
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}
	for _, stage := range stages {
		stageOut := stage(in)
		out := make(Bi, 1)
		go processStage(done, stageOut, out)
		in = out
	}
	return in
}
