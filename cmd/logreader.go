package cmd

import (
	"encoding/csv"
	"io"
)

type logReader struct {
	csvReader *csv.Reader
}

type logItem struct {
	row []string
	err error
}

func newLogReader(log io.Reader) *logReader {
	return &logReader{
		csvReader: csv.NewReader(log),
	}

}
func (l *logReader) rows(out chan<- logItem) {
	_, err := l.csvReader.Read()
	if err != nil {
		out <- logItem{row: nil, err: err}
	} else {
		for {
			row, err := l.csvReader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				out <- logItem{row: nil, err: err}
				break
			}
			out <- logItem{row: row, err: nil}
		}
	}
	close(out)
}
