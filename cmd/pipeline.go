package cmd

// Pipeline represents a sequential pipeline of stages
// where processing can occur.  They must be configured
// with a data source before use.
type Pipeline struct {
	stages []Stage
}

// NewPipeline creates a new pipeline with the specific stages in sequential
// order.  If outputs is specified then those are attached as outputs to each
// respective stage.
func NewPipeline(s ...Stage) *Pipeline {
	return &Pipeline{stages: s}
}

// Start sets up the pipeline and starts it using the src as the source of the payload.
// The function returns a read channel from the final stage of the pipeline and
// an error read channel that can be used to retrieve errors from any stage.
// If outputs is specified then those channels will be used as an output for
// the corresponding stage and they will be appended to the end of the list of output channels for
// each stage.
func (p *Pipeline) Start(src Source, outputs ...chan<- Payload) (<-chan Payload, <-chan error) {

	stagesCh := make([]chan Payload, len(p.stages)+1)
	errCh := make(chan error, len(p.stages)+2)

	for i := 0; i < len(stagesCh); i++ {
		stagesCh[i] = make(chan Payload)

	}

	for i := 0; i < len(p.stages); i++ {

		// append any specified outputs to the end of the outputs list
		var outCh []chan<- Payload
		if i < len(outputs) && len(outputs[i:]) > 0 {
			outCh = []chan<- Payload{stagesCh[i+1], outputs[i]}
		} else {
			outCh = []chan<- Payload{stagesCh[i+1]}
		}

		go func(n int, outCh []chan<- Payload) {
			p.stages[n].Run(&logMonitorStageParams{stageNum: n, inChan: stagesCh[n], outChan: outCh, errChan: errCh})

			//Each goroutine is responsible for closing the downstream channel to signal that it is done.
			close(stagesCh[n+1])

		}(i, outCh)
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
