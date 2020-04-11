package cmd

type Payload interface {
	Clone() Payload
	IteratorKey() (int, error)
}

type Transformer interface {
	Transform([]Payload) Payload
}

type StageParams interface {
	Input() <-chan Payload
	Output() chan<- Payload
	Error() chan<- error
}

type SourceParams interface {
	Output() chan<- Payload
	Error() chan<- error
}

type Source interface {
	Data(SourceParams)
}

type Stage interface {
	Run(StageParams)
}

type Sink interface {
	Write()
}
