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

// CSVLogSource implements csv file data source
// for a pipeline
type CSVLogSource struct {
	csvReader *csv.Reader
}

// NewCSVLogSource takes a source reader and returns
// a CSV reader.
func NewCSVLogSource(log io.Reader) *CSVLogSource {
	return &CSVLogSource{
		csvReader: csv.NewReader(log),
	}
}

type csvLogSourceParams struct {
	outChan chan<- Payload
	errChan chan<- error
}

func (p *csvLogSourceParams) Output() chan<- Payload { return p.outChan }
func (p *csvLogSourceParams) Error() chan<- error    { return p.errChan }

// Data loops through the CSV designed by the
// CSVLogSource and sends it through the output channel
// designated by the SourceParams
func (l *CSVLogSource) Data(s SourceParams) {

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
