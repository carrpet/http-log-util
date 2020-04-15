package cmd

// this code will start the Pipeline

// pipeline should contain a go channel that listens on a channel and
// consumes 10 seconds worth of data

type pipeline struct {
	stages []Stage
}

func newPipeline(s ...Stage) *pipeline {
	return &pipeline{stages: s}
}

// Start sets up the pipeline and starts it using the src as the source of the payload.
// The function returns a read channel from the final stage of the pipeline and
// an error read channel that can be used to retrieve errors from any stage.
func (p *pipeline) Start(src Source) (<-chan Payload, <-chan error) {

	stagesCh := make([]chan Payload, len(p.stages)+1)
	errCh := make(chan error, len(p.stages)+2)

	for i := 0; i < len(stagesCh); i++ {
		stagesCh[i] = make(chan Payload)

	}
	for i := 0; i < len(p.stages); i++ {
		go func(n int) {
			p.stages[n].Run(&logMonitorStageParams{stageNum: n, inChan: stagesCh[n], outChan: stagesCh[n+1], errChan: errCh})

			//Each goroutine is responsible for closing the downstream channel to signal that it is done.
			close(stagesCh[n+1])

		}(i)
	}

	// start source goroutine
	go func() {

		// There is only one SourceParam implementation here so we can use it
		// but ideally we would have this abstracted to be source agnostic.
		src.Data(&csvLogSourceParams{outChan: stagesCh[0], errChan: errCh})
		close(stagesCh[0])
	}()

	return stagesCh[len(stagesCh)-1], errCh

}
