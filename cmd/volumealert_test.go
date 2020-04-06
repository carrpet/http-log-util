package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

type volumeAlertTest struct {
	input []requestVolume
	want  []string
}

var (
	reqs     = []int{6, 2, 10, 200, 400, 3, 9, 1}
	endTimes = []int{1549573863, 1549573865, 1549573867, 1549573870, 1549573872,
		1549573874, 1549573876, 1549573878}
	want = []string{
		fmt.Sprintf(highTrafficAlertMsg, 210),
		fmt.Sprintf(alertRecoverMsg, time.Unix(1549573878, 0).Format(time.RFC3339)),
	}
)

func setUpTestDataNoError() []requestVolume {
	result := []requestVolume{}
	for i := range reqs {
		result = append(result, *newRequestVolume(reqs[i], nil, endTimes[i]))
	}
	return result

}

func TestVolumeAlert(t *testing.T) {

	// every entry represents two seconds, need 12 reqs per second
	// evaluated every 4 seconds
	cfg := volumeAlertConfig{alertThreshold: 12, alertFrequency: 4, interval: 2}
	testData := setUpTestDataNoError()
	volChan := make(chan requestVolume)
	go func() {
		for _, x := range testData {
			volChan <- x
		}
		close(volChan)
	}()
	results := []byte{}
	writer := bytes.NewBuffer(results)
	cfg.requestVolumeAlert(volChan, writer)
	expected := strings.Join(want, "")
	actual := writer.String()
	if actual != expected {
		t.Errorf(testErrMessage("Log contents differed", expected, actual))
	}

}
