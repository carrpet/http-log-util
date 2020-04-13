package cmd

import (
	"fmt"
	"strconv"
	"testing"
)

var (
	testErrMessage = func(msg, expected, actual string) string {
		return fmt.Sprintf(msg+": "+"expected: %s, actual: %s", expected, actual)
	}
)

type mockStage struct {
}

func (m mockStage) Run(p StageParams) {

	for data := range p.Input() {
		p.Output() <- data
	}
}

type mockPayload struct {
	data int
}

func (m mockPayload) IteratorKey() (int, error) {
	return m.data, nil
}

type mockSource struct {
	data []int
}

func (m mockSource) Data(s SourceParams) {

	for _, val := range m.data {
		s.Output() <- mockPayload{data: val}
	}
	close(s.Output())
	close(s.Error())

}

type mockSink struct {
	expected []int
	t        *testing.T
}

func (m mockSink) Write(params SinkParams) {
	//assert that the incoming data is the data sent through the pipeline
	length := len(m.expected)
	i := 0
	for x := range params.Input() {

		if length < 0 {
			m.t.Errorf("Actual length greater expected number of payload items")
		}
		output, ok := x.(mockPayload)
		if !ok {
			m.t.Errorf("Expected and did not receive a mockPayload object")
		}
		if output.data != m.expected[i] {
			m.t.Errorf(testErrMessage("Pipeline output validation failed", strconv.Itoa(i), strconv.Itoa(output.data)))
		}
		i++
		length = length - 1
	}
	if length > 0 {
		m.t.Errorf("Actual length less then expected number of payload items")
	}

}

func TestPipelineStart(t *testing.T) {

	toTest := newPipeline(mockStage{}, mockStage{})
	testPayload := []int{1, 2, 3, 4, 5, 6}
	src := mockSource{data: testPayload}
	sink := mockSink{expected: testPayload, t: t}
	toTest.Start(src, sink)

}
