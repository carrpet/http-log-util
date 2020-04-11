package cmd

type StageConfig struct {
	interval int
	proc     Transformer
}

func (s StageConfig) Run(params StageParams) {

}

type LogMonitorStageParams struct {
	stageNum int
	inChan   <-chan Payload
	outChan  chan<- Payload
	errChan  chan<- error
}

func (s *LogMonitorStageParams) Input() <-chan Payload { return s.inChan }

func (s *LogMonitorStageParams) Output() chan<- Payload { return s.outChan }

func (s *LogMonitorStageParams) Error() chan<- error { return s.errChan }

func (s *LogMonitorStageParams) Run(params StageParams) {

}
