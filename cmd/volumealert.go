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

func (s HttpStats) Error() bool {
	return false
}

type volumeAlertConfig struct {
	alertThreshold int
	alertFrequency int
	interval       int
}

//TODO
func newRequestVolume(numReq int, err error, ts int) *requestVolume {
	return &requestVolume{numRequests: numReq, ts: timestamp{startTime: ts}}
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
					time.Unix(int64(x.ts.endTime), 0).Format(time.RFC3339))
			}
			totalRequests = 0
		}
	}

}
