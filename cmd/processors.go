package cmd

import (
	"strconv"
	"time"
)

type logItems []logItem

type logItem struct {
	row []string
}

func getLogItemTime(l *logItem) (int, error) {
	return strconv.Atoi(l.row[date])
}

func (l *logItem) StartTime() int {
	result, _ := getLogItemTime(l)
	return result
}

func (l *logItem) EndTime() int {
	result, _ := getLogItemTime(l)
	return result
}

type timestamp struct {
	startTime int
	endTime   int
}

type requestVolume struct {
	numRequests int
	ts          timestamp
}

func (r *requestVolume) StartTime() int { return r.ts.startTime }

func (r *requestVolume) EndTime() int { return r.ts.endTime }

// represents a state transition for volume alerts
type volumeAlertTransition struct {
	alertFired bool
	timestamp  time.Time
}

type volumeAlerts struct {
	alerts []volumeAlertTransition
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
	transformFunc func(logItems, timestamp) Payload
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
	ts := timestamp{startTime: p[0].StartTime(), endTime: p[len(p)-1].EndTime()}
	payload := convertToLogItems(p)
	return r.transformFunc(payload, ts)
}

func requestVolumeTransformFunc(payload logItems, ts timestamp) Payload {
	return &requestVolume{numRequests: len(payload), ts: ts}

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
