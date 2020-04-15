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

func (hs httpStats) Write(w io.Writer) {
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Top Hits For Time Period %s to %s:\n",
		time.Unix(int64(hs.StartTime()), 0).Format(time.RFC3339Nano),
		time.Unix(int64(hs.EndTime()), 0).Format(time.RFC3339Nano))
	for _, th := range hs.topHits {
		fmt.Fprintf(w, "Section: %s, Number of Hits: %s\n", th.section, th.hits)
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Non Success Status Requests For Time Period %s to %s:\n",
		time.Unix(int64(hs.StartTime()), 0).Format(time.RFC3339Nano),
		time.Unix(int64(hs.EndTime()), 0).Format(time.RFC3339Nano))
	for _, us := range hs.unsuccessfulReqs {
		fmt.Fprintf(w, "Section: %s, Number of Non-Success Requests: %d\n",
			us.section, us.hits)
	}
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
