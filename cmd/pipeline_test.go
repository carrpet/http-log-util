package cmd

import (
	"fmt"
	"strconv"
	"testing"
)

var testErrMessage = func(msg, expected, actual string) string {
	return fmt.Sprintf(msg+": "+"expected: %s, actual: %s", expected, actual)
}

type mockStage struct {
	toAdd int
}

func (m mockStage) Run(p StageParams) {

	for data := range p.Input() {
		intdata := data.(mockPayload)
		intdata.data = intdata.data + m.toAdd
		p.Output()[0] <- intdata
	}
}

type mockPayload struct {
	data int
}

func (m mockPayload) StartTime() int { return m.data }

func (m mockPayload) EndTime() int { return m.data }

type mockSource struct {
	data []int
}

func (m mockSource) Data(s SourceParams) {

	for _, val := range m.data {
		s.Output() <- mockPayload{data: val}
	}

}

func TestPipelineStart(t *testing.T) {

	// create a pipeline that adds 2 and 8 to each input respectively
	toTest := NewPipeline(mockStage{toAdd: 2}, mockStage{toAdd: 8})
	testPayload := []int{1, 2, 3, 4, 5, 6}
	expected := []int{11, 12, 13, 14, 15, 16}

	src := mockSource{data: testPayload}
	sinkCh, _ := toTest.Start(src)

	//assert that the incoming data is the data sent through the pipeline
	length := len(expected)
	i := 0
	for x := range sinkCh {

		if length < 0 {
			t.Errorf("Actual length greater expected number of payload items")
		}
		output, ok := x.(mockPayload)
		if !ok {
			t.Errorf("Expected and did not receive a mockPayload object")
		}
		if output.data != expected[i] {
			t.Errorf(testErrMessage("Pipeline output validation failed", strconv.Itoa(i), strconv.Itoa(output.data)))
		}
		i++
		length--
	}
	if length > 0 {
		t.Errorf("Actual length less then expected number of payload items")
	}

}
