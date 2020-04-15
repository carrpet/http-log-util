package cmd

import (
	"fmt"
	"strconv"
	"testing"
)

func TestHttpStatsTransformFunc(t *testing.T) {
	testdata := logItems{
		{row: []string{"10.0.0.2", "-", "apache", "1549573862", "GET /api/user/bleh/h HTTP/1.0", "200", "1234"}},
		{row: []string{"10.0.0.4", "-", "apache", "1549573861", "GET /api/user HTTP/1.0", "200", "1234"}},
		{row: []string{"10.0.0.1", "-", "apache", "1549573862", "GET /api/help HTTP/1.0", "500", "1136"}},
		{row: []string{"10.0.0.4", "-", "apache", "1549573862", "POST /api/help HTTP/1.0", "500", "1234"}},
		{row: []string{"10.0.0.1", "-", "apache", "1549573862", "GET /api/help HTTP/1.0", "200", "1234"}},
		{row: []string{"10.0.0.1", "-", "apache", "1549573862", "GET /report HTTP/1.0", "200", "1194"}},
	}

	toTest := NewHTTPStatsProcessor()
	result := toTest.transformFunc(testdata, timestamp{}).(*httpStats)
	topSection := result.topHits[0].section
	topHitsCount := result.topHits[0].hits
	expected := httpStats{
		topHits:          []topHitStat{{section: "/api", hits: "5"}},
		unsuccessfulReqs: []nonSuccessHits{{section: "/api", hits: 2}, {section: "/report", hits: 0}}}
	if topSection != expected.topHits[0].section {
		t.Errorf(
			testErrMessage("Returned section was incorrect",
				"section should be /api", fmt.Sprintf("section was: %s", topSection)))
	}
	if topHitsCount != expected.topHits[0].hits {
		t.Errorf(
			testErrMessage("Number of hits was incorrect", "hits should be 5",
				fmt.Sprintf("Hits was %s", topHitsCount)))
	}
	for _, ex := range expected.unsuccessfulReqs {

		// we need to loop through all the results because order is not guaranteed
		found := false
		for _, actual := range result.unsuccessfulReqs {
			if ex.section == actual.section {
				found = true
				if ex.hits != actual.hits {
					t.Errorf(
						testErrMessage("UnsuccessfulReqs hits was different than expected",
							fmt.Sprintf("%d", ex.hits), fmt.Sprintf("%d", actual.hits)))
				}
			}
		}
		if !found {
			t.Errorf(
				testErrMessage("Missing unsuccessfulReqs section from result", ex.section, ""))
		}
	}

}

func TestRequestVolumeProcessorTransformFunc(t *testing.T) {
	testdata := logItems{
		{row: []string{"10.0.0.2", "-", "apache", "1549573862", "GET /api/user/bleh/h HTTP/1.0", "200", "1234"}},
		{row: []string{"10.0.0.4", "-", "apache", "1549573861", "GET /api/user HTTP/1.0", "200", "1234"}},
		{row: []string{"10.0.0.1", "-", "apache", "1549573862", "GET /api/help HTTP/1.0", "500", "1136"}},
		{row: []string{"10.0.0.4", "-", "apache", "1549573862", "POST /api/help HTTP/1.0", "200", "1234"}},
		{row: []string{"10.0.0.1", "-", "apache", "1549573862", "GET /api/help HTTP/1.0", "200", "1234"}},
		{row: []string{"10.0.0.1", "-", "apache", "1549573863", "GET /report HTTP/1.0", "500", "1194"}},
	}
	expectedTimeStart, _ := strconv.Atoi(testdata[0].row[date])
	expectedTimeEnd, _ := strconv.Atoi(testdata[5].row[date])
	timeStamp := timestamp{startTime: expectedTimeStart, endTime: expectedTimeEnd}
	expected := requestVolume{numRequests: 6, ts: timeStamp}
	rvProcessor := NewRequestVolumeProcessor()
	result := rvProcessor.transformFunc(testdata, timeStamp).(*requestVolume)
	if result.numRequests != expected.numRequests {
		t.Errorf(testErrMessage("RequestVolumeProcessor.transformFunc had wrong output",
			"numRequests == "+strconv.Itoa(expected.numRequests), "numRequests == "+strconv.Itoa(result.numRequests)))
	}
	if result.ts != expected.ts {
		t.Errorf(testErrMessage("RequestVolumeProcessor.transformFunc had wrong output",
			"startTime / endTime ==  "+strconv.Itoa(expected.ts.startTime)+"/"+strconv.Itoa(expected.ts.endTime),
			"startTime / endTime ==  "+strconv.Itoa(result.ts.startTime)+"/"+strconv.Itoa(result.ts.endTime)))
	}

}

func TestAlertOutputProcessorTransformFunc(t *testing.T) {

	// tell function to alert for request volume > 2 requests per sec
	toTest := NewAlertProcessor(2)

	data := [3]requestVolumes{
		// 109 requests in 31 seconds
		{{numRequests: 1, ts: timestamp{startTime: 1, endTime: 10}},
			{numRequests: 8, ts: timestamp{startTime: 11, endTime: 21}},
			{numRequests: 100, ts: timestamp{startTime: 22, endTime: 32}}},

		// 7 requests in 30 seconds
		{{numRequests: 1, ts: timestamp{startTime: 33, endTime: 43}},
			{numRequests: 4, ts: timestamp{startTime: 43, endTime: 54}},
			{numRequests: 2, ts: timestamp{startTime: 55, endTime: 63}}},

		// 3 requests in 30 seconds
		{{numRequests: 1, ts: timestamp{startTime: 64, endTime: 74}},
			{numRequests: 1, ts: timestamp{startTime: 75, endTime: 85}},
			{numRequests: 1, ts: timestamp{startTime: 85, endTime: 94}}}}

	expected := [3]volumeAlertStatus{
		{alertFiring: true, time: 32, volume: 109},
		{alertFiring: false, time: 63, volume: 7},
		{alertFiring: false, time: 94, volume: 3}}

	for i, val := range data {
		result := toTest.transformFunc(val).(*volumeAlertStatus)
		if result.alertFiring != expected[i].alertFiring {
			t.Errorf(testErrMessage("alertOutputTransformFunc, wrong values for alertFiring",
				strconv.FormatBool(expected[i].alertFiring), strconv.FormatBool(result.alertFiring)))
		}
		if result.time != expected[i].time {
			t.Errorf(testErrMessage("alertOutputTransformFunc, wrong values for time",
				strconv.Itoa(expected[i].time), strconv.Itoa(result.time)))
		}
		if result.volume != expected[i].volume {
			t.Errorf(testErrMessage("alertOutputTransformFunc, wrong values for volume",
				strconv.Itoa(expected[i].volume), strconv.Itoa(result.volume)))
		}
	}

}
