package cmd

import (
	"encoding/csv"
	"io"
)

type logReader struct {
	csvReader *csv.Reader
}

func newLogReader(log io.Reader) *logReader {
	return &logReader{
		csvReader: csv.NewReader(log),
	}

}
func (l *logReader) rows() <-chan []string {
	_, err := l.csvReader.Read()
	if err != nil {
		// TODO:do something with the errors other than returning
		return nil
	}
	out := make(chan []string)
	go func() {
		for {
			row, err := l.csvReader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				//TODO: do something to indicate to the caller
				// that some non EOF error happened
				break
			}
			out <- row
		}
		close(out)
	}()
	return out
}
