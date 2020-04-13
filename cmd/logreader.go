package cmd

import (
	"encoding/csv"
	"io"
)

// represents the fields of the log file
const (
	rh       = iota
	rfc931   = iota
	authuser = iota
	date     = iota
	req      = iota
	status   = iota
	numBytes = iota
)

//logReader implements Source interface
type csvLogSource struct {
	csvReader *csv.Reader
}

func newCsvLogSource(log io.Reader) *csvLogSource {
	return &csvLogSource{
		csvReader: csv.NewReader(log),
	}

}

type csvLogSourceParams struct {
	outChan chan<- Payload
	errChan chan<- error
}

func (p *csvLogSourceParams) Output() chan<- Payload { return p.outChan }

func (p *csvLogSourceParams) Error() chan<- error { return p.errChan }

type sinkParams struct {
	inChan <-chan Payload
}

func (p *sinkParams) Input() <-chan Payload { return p.inChan }

func (l *csvLogSource) Data(s SourceParams) {

	//expect to read the header
	_, err := l.csvReader.Read()
	if err != nil {
		s.Error() <- err
	} else {
		for {
			row, err := l.csvReader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				s.Error() <- err
				break
			}
			s.Output() <- &logItem{row: row}
		}
	}
}
