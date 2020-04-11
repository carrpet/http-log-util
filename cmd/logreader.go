package cmd

import (
	"encoding/csv"
	"io"
)

//logReader implements Source interface
type csvLogSource struct {
	csvReader *csv.Reader
}

type logItem struct {
	row []string
}

func (li logItem) Error() bool {
	return li.err != nil
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
