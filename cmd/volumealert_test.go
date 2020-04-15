package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestVolumeAlert(t *testing.T) {

	testData := []Payload{
		&volumeAlertStatus{alertFiring: false, time: 5, volume: 10},
		&volumeAlertStatus{alertFiring: true, time: 11, volume: 100},
		&volumeAlertStatus{alertFiring: true, time: 18, volume: 500},
		&volumeAlertStatus{alertFiring: false, time: 24, volume: 40},
	}

	inCh := make(chan Payload)
	go func() {
		for _, val := range testData {
			inCh <- val
		}
		close(inCh)
	}()
	resultsBuf := []byte{}
	writer := bytes.NewBuffer(resultsBuf)
	expected := strings.Join([]string{
		fmt.Sprintf(highTrafficAlertMsg, 100),
		fmt.Sprintf(alertRecoverMsg, time.Unix(24, 0).Format(time.RFC3339))}, "")
	writeAlerts(inCh, writer)
	actual := writer.String()
	if expected != actual {
		t.Errorf(testErrMessage("Alert writer results differed", expected, actual))
	}
}
