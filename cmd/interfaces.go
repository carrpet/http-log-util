package cmd

// Payload represents operations that can be performed
// on an pipeline payload object.
type Payload interface {
	StartTime() int
	EndTime() int
}

// Transformer represents actions that can
// be performed to input payload at each stage, sending
// the output downsteam.
type Transformer interface {
	Transform([]Payload) Payload
}

// StageParams represents configuration for the
// input, output, and error channels for each stage.
type StageParams interface {
	Input() <-chan Payload
	Output() []chan<- Payload
	Error() chan<- error
}

// SourceParams represents configuration for the
// output and error channels for a pipeline source.
type SourceParams interface {
	Output() chan<- Payload
	Error() chan<- error
}

// Source represents actions that can be performed
// by a pipeline's source.
type Source interface {
	Data(SourceParams)
}

// Stage represents actions that can be performed
// by each stage of the pipeline.
type Stage interface {
	Run(StageParams)
}

// SinkParams represents channel configuration
// for a pipeline's sink.
type SinkParams interface {
	Input() <-chan Payload
}

// Sink represents actions that can be performed
// by a pipeline's sink.
type Sink interface {
	Write(SinkParams)
}
