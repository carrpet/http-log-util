package cmd

import (
	"strconv"
	"testing"
	"time"
)

/*
func TestComputeTopHitsReturnsCorrectData(t *testing.T) {
	testdata := [][]string{
		[]string{"10.0.0.2", "-", "apache", "1549573862", "GET /api/user/bleh/h HTTP/1.0", "200", "1234"},
		[]string{"10.0.0.4", "-", "apache", "1549573861", "GET /api/user HTTP/1.0", "200", "1234"},
		[]string{"10.0.0.1", "-", "apache", "1549573862", "GET /api/help HTTP/1.0", "500", "1136"},
		[]string{"10.0.0.4", "-", "apache", "1549573862", "POST /api/help HTTP/1.0", "200", "1234"},
		[]string{"10.0.0.1", "-", "apache", "1549573862", "GET /api/help HTTP/1.0", "200", "1234"},
		[]string{"10.0.0.1", "-", "apache", "1549573862", "GET /report HTTP/1.0", "500", "1194"},
	}
	result := computeTopHits(testdata)
	topSection := result.topHits[0].section
	topHitsCount := result.topHits[0].hits
	if topSection != "/api" {
		t.Errorf(
			testErrMessage("Returned section was incorrect",
				"section should be /api", fmt.Sprintf("section was: %s", topSection)))
	}
	if topHitsCount != "5" {
		t.Errorf(
			testErrMessage("Number of hits was incorrect", "hits should be 5",
				fmt.Sprintf("Hits was %s", topHitsCount)))
	}

}
*/

func TestRequestVolumeProcessorTransformFunc(t *testing.T) {
	testdata := logItems{
		{row: []string{"10.0.0.2", "-", "apache", "1549573862", "GET /api/user/bleh/h HTTP/1.0", "200", "1234"}},
		{row: []string{"10.0.0.4", "-", "apache", "1549573861", "GET /api/user HTTP/1.0", "200", "1234"}},
		{row: []string{"10.0.0.1", "-", "apache", "1549573862", "GET /api/help HTTP/1.0", "500", "1136"}},
		{row: []string{"10.0.0.4", "-", "apache", "1549573862", "POST /api/help HTTP/1.0", "200", "1234"}},
		{row: []string{"10.0.0.1", "-", "apache", "1549573862", "GET /api/help HTTP/1.0", "200", "1234"}},
		{row: []string{"10.0.0.1", "-", "apache", "1549573863", "GET /report HTTP/1.0", "500", "1194"}},
	}
	expectedTimeStr, _ := strconv.Atoi(testdata[5].row[date])
	expected := requestVolume{numRequests: 6, endTime: time.Unix(int64(expectedTimeStr), 0)}
	rvProcessor := newRequestVolumeProcessor()
	result := rvProcessor.transformFunc(testdata).(*requestVolume)
	if result.numRequests != expected.numRequests {
		t.Errorf(testErrMessage("RequestVolumeProcessor.transformFunc had wrong output",
			"numRequests == "+strconv.Itoa(expected.numRequests), "numRequests == "+strconv.Itoa(result.numRequests)))
	}
	if result.endTime != expected.endTime {
		t.Errorf(testErrMessage("RequestVolumeProcessor.transformFunc had wrong output",
			"endTime == "+expected.endTime.String(), "endTime == "+result.endTime.String()))
	}

}

/*
func TestAlertOutputProcessorTransformFunc(t *testing.T) {

	// tell function to alert for request volume > 2 requests per sec
	// with function expecting 10 seconds of data
	toTest := newAlertOutputProcessor(2)
	testFunc := toTest.transformFunc(10)
	testdata := []requestVolume{
		{numRequests: 1, endTime: time.Unix(10,0)},
		{numRequests: 8, endTime: time.Unix(12,0)},
		{numRequests: 10, endTime: time.Unix()}
	}


}
*/
