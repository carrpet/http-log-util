package cmd

import (
	"fmt"
	"strings"
	"testing"
)

var (
	readerData = `"remotehost","rfc931","authuser","date","request","status","bytes"
"10.0.0.2","-","apache",1549573860,"GET /api/user HTTP/1.0",200,1234
"10.0.0.4","-","apache",1549573861,"GET /api/user HTTP/1.0",200,1136
"10.0.0.5","-","apache",1549573861,"GET /api/user HTTP/1.0",200,1194
"10.0.0.1","-","apache",1549573861,"GET /api/user HTTP/1.0",200,1261`

	invalidFile = `"remotehost","rfc931","authuser","date","request","status","bytes"
"10.0.0.2","-","apache",1549573860,"GET /api/user HTTP/1.0",200,1234
"10.0.0.4","-","a`
)

func chanHelper(data string) *logReader {
	return newLogReader(strings.NewReader(data))
}

func TestLogReaderProducesAllRowsNoHeader(t *testing.T) {
	lReader := chanHelper(readerData)
	lChan := make(chan Iterable)

	// read in the rows to channel
	go lReader.rows(lChan)
	counter := 0
	for range lChan {
		counter++
	}

	if counter != 4 {
		t.Error(testErrMessage("LogReader produced incorrect number of rows",
			"counter is 4",
			fmt.Sprintf("counter is %d", counter)))
	}
}

func TestLogReaderHandlesInvalidFile(t *testing.T) {
	lReader := chanHelper(invalidFile)
	lChan := make(chan Iterable)
	go lReader.rows(lChan)
	counter := 0
	for range lChan {
		counter++
	}
	if counter != 2 {
		t.Error("Counter was wrong!")
	}

}
