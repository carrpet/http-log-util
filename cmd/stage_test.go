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

func TestLogMonitorStageSendsIntervalLengthData(t *testing.T) {
	data := []int{0, 2, 3, 4, 5, 6, 7, 10}
	expected := []int{5, 22, 10}
	toTest := newStage(mockTransformer{}, 3)
	inCh := make(chan Payload)
	outCh := make(chan Payload)
	//errCh := make(chan error)
	go func() {
		for _, val := range data {
			inCh <- mockPayload{data: val}
		}
		close(inCh)

	}()

	go func() {
		expectedLen := len(expected)
		i := 0
		for x := range outCh {
			val := x.(mockPayload)
			if expectedLen < 0 {
				t.Errorf("Output length was greater than expected length")
			}
			if val.data != expected[i] {
				t.Errorf(testErrMessage("Expected value for stage was not equal to actual stage output",
					strconv.Itoa(expected[i]), strconv.Itoa(val.data)))
			}
			i++
			expectedLen--
		}
		close(outCh)

	}()
	params := &LogMonitorStageParams{stageNum: 0, inChan: inCh, outChan: outCh, errChan: nil}
	toTest.Run(params)
}
