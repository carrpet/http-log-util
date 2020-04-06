package cmd

import (
	"fmt"
	"io"
	"time"
)

/*
type DataSource interface
	Read() interface{}


type ChanDataSource <-chan interface{}

func (ch *ChanDataSource) Read() interface{} {
	return <- ch
}

// Console Logger represents a sink of data which can write to an output
type ConsoleLogger interface
	WriteLog(io.Writer)
	RegisterSource()

*/
// requestVolume represents the number of requests
// processed over the time interval (in seconds)
// and err will be propagated if an error happened
// upstream

const (
	highTrafficAlertMsg = "High traffic generated an alert - hits = %d"
	alertRecoverMsg     = "The alert has recovered at time %s"
)

type requestVolume struct {
	numRequests int
	err         error
	endTime     time.Time //indicates the end time for this interval
}

type volumeAlertConfig struct {
	alertThreshold int
	alertFrequency int
	interval       int
}

func newRequestVolume(numReq int, err error, ts int) *requestVolume {
	return &requestVolume{numRequests: numReq, err: err, endTime: time.Unix(int64(ts), 0)}
}

func (c *volumeAlertConfig) requestVolumeAlert(ch <-chan requestVolume, w io.Writer) {

	logDuration := 0
	totalRequests := 0
	highState := false

	totalRequestThreshold := c.alertThreshold * c.alertFrequency
	for x := range ch {
		println("Reading from channel!")
		logDuration = logDuration + c.interval
		totalRequests = totalRequests + x.numRequests
		if logDuration%c.alertFrequency == 0 {
			//trigger alert
			if !highState && totalRequests >= totalRequestThreshold {
				highState = true
				fmt.Fprintf(w, highTrafficAlertMsg, totalRequests)
			} else if highState && totalRequests < totalRequestThreshold {
				highState = false
				fmt.Fprintf(w, alertRecoverMsg,
					x.endTime.Format(time.RFC3339))
			}
			totalRequests = 0
		}
	}

}
