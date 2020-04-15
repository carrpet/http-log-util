package cmd

import (
	"strings"
	"testing"
)

var (
	readerData = `"remotehost","rfc931","authuser","date","request","status","bytes"
"10.0.0.2","-","apache",1549573860,"GET /api/user HTTP/1.0",200,1234
"10.0.0.4","-","apache",1549573861,"GET /api/user HTTP/1.0",200,1136
"10.0.0.5","-","apache",1549573861,"GET /api/user HTTP/1.0",200,1194
"10.0.0.1","-","apache",1549573861,"GET /api/user HTTP/1.0",200,1261`
)

func TestLogReaderProducesLogItemRowsTimesNoHeader(t *testing.T) {

	src := NewCSVLogSource(strings.NewReader(readerData))
	ch := make(chan Payload)
	params := &csvLogSourceParams{outChan: ch}

	go func() {
		src.Data(params)
		close(ch)
	}()
	expectedTimes := []string{"1549573860", "1549573861", "1549573861", "1549573861"}
	i := 0
	for x := range ch {
		li := x.(*logItem)
		if li.row[date] != expectedTimes[i] {
			t.Error(testErrMessage("LogReader produced incorrect expected time",
				expectedTimes[i],
				li.row[date]))
		}
		i++
	}
}
