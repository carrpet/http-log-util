package cmd

type StageConfig struct {
	interval int
	proc     Transformer
}

func newStage(proc Transformer, interval int) Stage {
	return StageConfig{proc: proc, interval: interval}
}

func (s StageConfig) Run(p StageParams) {
	lf := LogFilter{interval: s.interval}
	processFunc := lf.ProcessInterval()
	var buf []Payload
	for x := range p.Input() {
		buf = append(buf, x)
		//TODO: handle error
		it, err := x.IteratorKey()
		if err != nil {
			p.Error() <- err
			return
		}
		iVal := processFunc(it)
		if !iVal {
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
	//TODO: maybe close the channel here?
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
