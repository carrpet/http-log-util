package cmd

type StageConfig struct {
	interval int
	proc     Transformer
}

func newStage(proc Transformer, interval int) Stage {
	return StageConfig{proc: proc, interval: interval}
}

func (s StageConfig) Run(p StageParams) {

	findInterval := processInterval(s.interval)
	var buf []Payload

	for x := range p.Input() {
		buf = append(buf, x)
		recordInterval := findInterval(x)
		if !recordInterval {
			continue
		}
		result := s.proc.Transform(buf)
		p.Output() <- result
		buf = nil
	}

	//handle the case where we are at the end of input and need to send
	// the remaining items
	if buf != nil {
		p.Output() <- s.proc.Transform(buf)
	}

}

// LogMonitorStageParams specifies the input, output, and error channels for the stage.
type LogMonitorStageParams struct {
	stageNum int
	inChan   <-chan Payload
	outChan  chan<- Payload
	errChan  chan<- error
}

func (s *LogMonitorStageParams) Input() <-chan Payload { return s.inChan }

func (s *LogMonitorStageParams) Output() chan<- Payload { return s.outChan }

func (s *LogMonitorStageParams) Error() chan<- error { return s.errChan }
