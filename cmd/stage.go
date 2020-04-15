package cmd

type stageConfig struct {
	interval int
	proc     []Transformer
}

// newStage creates a new stage with one output and associated transformer
func newStage(proc Transformer, interval int) Stage {
	return stageConfig{proc: []Transformer{proc}, interval: interval}
}

// NewFanOutStage allows multiple Transformers to process
// data to be sent to multiple downstream stages. The first
// transformer specified will be used for the first output channel
// the second transformer will be used for the second...
// If stage output channels are specified in the Pipeline.Start method, then
// they will correspond to the last transformers specified in the list
func NewFanOutStage(procs []Transformer, interval int) Stage {
	return stageConfig{proc: procs, interval: interval}
}

func (s stageConfig) Run(p StageParams) {

	findInterval := processInterval(s.interval)
	var buf []Payload

	for x := range p.Input() {
		buf = append(buf, x)
		recordInterval := findInterval(x)
		if !recordInterval {
			continue
		}
		for i := range p.Output() {
			result := s.proc[i].Transform(buf)
			p.Output()[i] <- result
		}
		buf = nil
	}

	//handle the case where we are at the end of input and need to send
	// the remaining items
	if buf != nil {
		for i := range p.Output() {
			result := s.proc[i].Transform(buf)
			p.Output()[i] <- result
		}
	}

}

// logMonitorStageParams specifies the input, output, and error channels for the stage.
type logMonitorStageParams struct {
	stageNum int
	inChan   <-chan Payload
	outChan  []chan<- Payload
	errChan  chan<- error
}

func (s *logMonitorStageParams) Input() <-chan Payload { return s.inChan }

func (s *logMonitorStageParams) Output() []chan<- Payload { return s.outChan }

func (s *logMonitorStageParams) Error() chan<- error { return s.errChan }
