package cmd

import (
	"fmt"
	"io"
	"time"
)

const (
	highTrafficAlertMsg = "High traffic generated an alert - hits = %d"
	alertRecoverMsg     = "The alert has recovered at time %s"
)

type requestVolume struct {
	numRequests int
	err         error
	endTime     time.Time //indicates the end time for this interval
}

func (rv requestVolume) IteratorKey() (int, error) {
	return 0, nil
}

func (rv requestVolume) Error() bool {
	return false

}

func (s HttpStats) Error() bool {
	return false
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
