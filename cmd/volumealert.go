package cmd

import (
	"fmt"
	"io"
	"time"
)

const (
	highTrafficAlertMsg = "High traffic generated an alert - hits = %d\n"
	alertRecoverMsg     = "The alert has recovered at time %s\n"
)

func newRequestVolume(numReq int, err error, ts int) *requestVolume {
	return &requestVolume{numRequests: numReq, ts: timestamp{startTime: ts}}
}

func writeAlerts(ch <-chan Payload, w io.Writer) {

	alertFiring := false
	for x := range ch {
		alertStatus := x.(*volumeAlertStatus)
		if !alertFiring && alertStatus.alertFiring {
			alertFiring = true
			fmt.Fprintf(w, highTrafficAlertMsg, alertStatus.volume)
		} else if alertFiring && !alertStatus.alertFiring {
			alertFiring = false
			fmt.Fprintf(w, alertRecoverMsg,
				time.Unix(int64(x.EndTime()), 0).Format(time.RFC3339))
		}
	}
}
