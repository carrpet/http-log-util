package cmd

// this code will start the Pipeline

// pipeline should contain a go channel that listens on a channel and
// consumes 10 seconds worth of data

type Pipeline struct {
	stages []Stage
}

func newPipeline(s ...Stage) *Pipeline {
	return &Pipeline{stages: s}
}

func (p *Pipeline) Start(src Source, sink Sink) {

	stagesCh := make([]chan Payload, len(p.stages)+1)
	errCh := make(chan error, len(p.stages)+2)

	for i := range p.stages {
		stagesCh[i] = make(chan Payload)

	}
	for i, val := range p.stages {
		go func() {
			val.Run(&LogMonitorStageParams{stageNum: i, inChan: stagesCh[i], outChan: stagesCh[i+1], errChan: errCh})

			//Each goroutine is responsible for closing the downstream channel to signal that it is done.
			close(stagesCh[i+1])

		}()
	}

	go func() {

		// There is only one SourceParam implementation here so we can use it
		// but ideally we would have this abstracted to be source agnostic.
		src.Data(&csvLogSourceParams{outChan: stagesCh[0], errChan: errCh})
		close(errCh)
	}()
}

func newStage(proc Transformer, interval int) Stage {
	return StageConfig{proc: proc, interval: interval}

}
