package cmd

import (
	"fmt"
	"strconv"
	"time"
)

type logItems []logItem

type logItem struct {
	row []string
}

type requestVolume struct {
	numRequests int
	err         error
	endTime     time.Time //indicates the end time for this interval
}

// represents a state transition for volume alerts
type volumeAlertTransition struct {
	alertFired bool
	timestamp  time.Time
}

type volumeAlerts struct {
	alerts []volumeAlertTransition
}

func (rv requestVolume) IteratorKey() (int, error) {
	return 0, nil
}

func (li *logItem) IteratorKey() (int, error) {
	ts, err := strconv.Atoi(li.row[date])
	if err != nil {
		// TODO
		fmt.Errorf("Error retrieving iterator key for logItem: %s", err.Error())
	}
	return ts, nil

}

//TransformForWrite computes the metrics for
// writing to output.
/*
func (li LogItems) TransformForWrite() Writable {
	hits := map[string]int{}
	for _, val := range li {
		req := val.row[req]
		path := strings.Split(req, " ")
		section := "/" + strings.SplitN(path[1], "/", 3)[1]
		hits[section]++
	}

	//find the max hits section
	maxHits := 0
	var maxSection string
	for sect, h := range hits {
		if h > maxHits {
			maxHits = h
			maxSection = sect
		}
	}
	return HttpStats{topHits: []TopHitStat{{section: maxSection, hits: strconv.Itoa(maxHits)}}}
}
*/

type httpStatsProcessor struct {
}

// Transform converts the data to the form that is expected
// in the downstream stage of the pipeline.
func (hs *httpStatsProcessor) Transform(p []Payload) Payload {
	return nil
}

func newHTTPStatsProcessor() *httpStatsProcessor {
	return &httpStatsProcessor{}

}

/*
func newAlertOutputProcessor(threshold, freq int) *alertOutputProcessor {
	return &alertOutputProcessor{alertThreshold: threshold, alertFrequency: freq}

}
*/
type requestVolumeProcessor struct {
	transformFunc func(logItems) Payload
}

func newRequestVolumeProcessor() *requestVolumeProcessor {
	return &requestVolumeProcessor{transformFunc: requestVolumeTransformFunc}
}

type requestVolumeFunc func(int) *requestVolumeProcessor

func (r requestVolumeFunc) Transform(p []Payload) Payload {

	return nil
}

// Transform counts the number of total requests in the interval.
func (r *requestVolumeProcessor) Transform(p []Payload) Payload {
	payload := convertToLogItems(p)
	return r.transformFunc(payload)
}

func requestVolumeTransformFunc(payload logItems) Payload {
	endTime, err := strconv.Atoi(payload[len(payload)-1].row[date])
	if err != nil {
		return &requestVolume{err: err}
	}
	return &requestVolume{numRequests: len(payload),
		err: nil, endTime: time.Unix(int64(endTime), 0)}

}

type alertOutputProcessor struct {
	alertThreshold int
	transformFunc  func(int) func([]requestVolume) []volumeAlerts
}

func newAlertOutputProcessor(threshold int) *alertOutputProcessor {
	return &alertOutputProcessor{transformFunc: alertOutputTransformFuncBuilder, alertThreshold: threshold}
}

func (a *alertOutputProcessor) Transform(p []Payload) Payload {
	return nil
}

func alertOutputTransformFuncBuilder(interval int) func([]requestVolume) []volumeAlerts {
	return func(payload []requestVolume) []volumeAlerts {
		return []volumeAlerts{}

	}
}

func convertToLogItems(p []Payload) logItems {
	var payload logItems
	for _, val := range p {
		payload = append(payload, *(val.(*logItem)))
	}
	return payload
}
