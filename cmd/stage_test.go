package cmd

import (
	"strconv"
	"testing"
)

type mockTransformer struct {
}

func (m mockTransformer) Transform(data []Payload) Payload {
	sum := 0
	for _, val := range data {
		intval := val.(mockPayload)
		sum += intval.data
	}
	return mockPayload{data: sum}

}

func stageTestRunner(toTest Stage, validatorFunc func(Payload, Payload),
	input []Payload, expected []Payload, t *testing.T) {

	inCh := make(chan Payload)
	outCh := make(chan Payload) //errCh := make(chan error)
	go func() {
		for _, val := range input {
			inCh <- val
		}
		close(inCh)

	}()

	params := &logMonitorStageParams{stageNum: 0, inChan: inCh, outChan: outCh, errChan: nil}
	go func() {
		toTest.Run(params)
		close(outCh)
	}()

	i := -1
	for x := range outCh {
		i++
		if i > len(expected)-1 {
			t.Errorf("Output length was greater than expected length")
		}
		validatorFunc(expected[i], x)
	}
	if i < len(expected)-1 {
		t.Errorf(testErrMessage("Output length less than expected length", strconv.Itoa(len(expected)),
			strconv.Itoa(i+1)))
	}

}

func TestLogMonitorStageSendsIntervalLengthData(t *testing.T) {
	data := []Payload{
		mockPayload{data: 0},
		mockPayload{data: 2},
		mockPayload{data: 3},
		mockPayload{data: 4},
		mockPayload{data: 5},
		mockPayload{data: 6},
		mockPayload{data: 7},
		mockPayload{data: 10}}
	expected := []Payload{mockPayload{data: 5}, mockPayload{data: 22}, mockPayload{data: 10}}
	toTest := newStage(mockTransformer{}, 3)
	validate := func(expected Payload, actual Payload) {

		result := actual.(mockPayload)
		ex := expected.(mockPayload).data
		if result.data != ex {
			t.Errorf(testErrMessage("Expected value for stage was not equal to actual stage output",
				strconv.Itoa(ex), strconv.Itoa(result.data)))
		}
	}
	stageTestRunner(toTest, validate, data, expected, t)
}

func TestRequestVolumeStage(t *testing.T) {
	//Integration Stage Tests
	data := []Payload{
		&logItem{row: []string{"", "", "", "12345"}},
		&logItem{row: []string{"", "", "", "12344"}},
		&logItem{row: []string{"", "", "", "12345"}},
		&logItem{row: []string{"", "", "", "12346"}},
		&logItem{row: []string{"", "", "", "12347"}},
		&logItem{row: []string{"", "", "", "12348"}},
		&logItem{row: []string{"", "", "", "12348"}},
		&logItem{row: []string{"", "", "", "12349"}},
		&logItem{row: []string{"", "", "", "12349"}}}
	expected := []Payload{
		&requestVolume{numRequests: 4, ts: timestamp{startTime: 12345, endTime: 12346}},
		&requestVolume{numRequests: 4, ts: timestamp{startTime: 12347, endTime: 12349}},
		&requestVolume{numRequests: 1, ts: timestamp{startTime: 12349, endTime: 12349}}}
	toTest := newStage(newRequestVolumeProcessor(), 2)

	validate := func(expected Payload, actual Payload) {
		result := actual.(*requestVolume)
		ex := expected.(*requestVolume)
		if result.numRequests != ex.numRequests {
			t.Errorf("Expected value for stage was not equal to actual stage output, expected: %d, actual: %d",
				ex.numRequests, result.numRequests)
		}

		if result.StartTime() != ex.StartTime() {
			t.Errorf(testErrMessage("Expected start time for stage was not equal to actual start time",
				strconv.Itoa(ex.StartTime()), strconv.Itoa(result.StartTime())))

		}
		if result.EndTime() != ex.EndTime() {
			t.Errorf(testErrMessage("Expected endTime for stage was not equal to actual end time",
				strconv.Itoa(ex.EndTime()), strconv.Itoa(result.EndTime())))

		}

	}
	stageTestRunner(toTest, validate, data, expected, t)

}
