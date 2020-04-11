package cmd

import (
	"strconv"
	"testing"
)

var testdata = LogItems{{}, {}, {}}

type mockFlowRegulator struct {
	interval int
}

func (m mockFlowRegulator) ProcessInterval() func(...int) bool {

	// return true after on every third argument
	counter := 0
	return func(num ...int) bool {
		for range num {
			counter++
		}
		if counter%m.interval == 0 {
			return true
		}
		return false
	}
}

type mockSendable struct {
	interval int
	t        *testing.T
}

type sequenceTestType int

func (s sequenceTestType) IteratorKey() (int, error) {
	return int(s), nil
}

func (m mockSendable) SendTransformation(items []Iterable) {
	if len(items) != m.interval {
		m.t.Errorf(testErrMessage("The length of the data sent was wrong", strconv.Itoa(m.interval), strconv.Itoa(len(items))))
	}
}

func LogItemsFilterTestAccumulatesDataUntilProcessIntervalSaysStop(t *testing.T) {
	var ints [6]int
	mockFR := mockFlowRegulator{interval: 3}
	logChan := make(chan Iterable)
	done := make(chan interface{})
	mockS := mockSendable{}
	logItemsFilter(mockFR, logChan, done, mockS)

	for i := range ints {
		logChan <- sequenceTestType(i)
	}
}
