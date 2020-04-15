package cmd

import (
	"strconv"
	"strings"
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

type httpStats struct {
	ts      timestamp
	topHits []topHitStat
}

func (s *httpStats) StartTime() int { return s.ts.startTime }
func (s *httpStats) EndTime() int   { return s.ts.endTime }

type topHitStat struct {
	section string
	hits    string
}

type requestVolume struct {
	numRequests int
	ts          timestamp
}

type requestVolumes []requestVolume

func (r *requestVolume) StartTime() int { return r.ts.startTime }
func (r *requestVolume) EndTime() int   { return r.ts.endTime }

// represents a state transition for volume alerts
type volumeAlertStatus struct {
	alertFiring bool
	time        int
	volume      int
}

func (v *volumeAlertStatus) StartTime() int { return v.time }
func (v *volumeAlertStatus) EndTime() int   { return v.time }

type httpStatsProcessor struct {
	transformFunc func(logItems, timestamp) Payload
}

// Transform converts the data to the form that is expected
// in the downstream stage of the pipeline.
func (hs *httpStatsProcessor) Transform(p []Payload) Payload {
	ts := getPayloadsTimestamp(p)
	payload := convertToLogItems(p)
	return hs.transformFunc(payload, ts)
}

func newHTTPStatsProcessor() *httpStatsProcessor {
	return &httpStatsProcessor{transformFunc: httpStatsTransformFunc}
}

func httpStatsTransformFunc(p logItems, ts timestamp) Payload {
	hits := map[string]int{}
	for _, val := range p {
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
	return &httpStats{ts: ts,
		topHits: []topHitStat{{section: maxSection, hits: strconv.Itoa(maxHits)}}}

}

type requestVolumeProcessor struct {
	transformFunc func(logItems, timestamp) Payload
}

func newRequestVolumeProcessor() *requestVolumeProcessor {
	return &requestVolumeProcessor{transformFunc: requestVolumeTransformFunc}
}

// Transform counts the number of total requests in the interval.
func (r *requestVolumeProcessor) Transform(p []Payload) Payload {
	ts := getPayloadsTimestamp(p)
	payload := convertToLogItems(p)
	return r.transformFunc(payload, ts)
}

func requestVolumeTransformFunc(payload logItems, ts timestamp) Payload {
	return &requestVolume{numRequests: len(payload), ts: ts}
}

type alertOutputProcessor struct {
	alertThreshold int
	transformFunc  func(requestVolumes) Payload
}

func newAlertOutputProcessor(threshold int) *alertOutputProcessor {
	return &alertOutputProcessor{transformFunc: alertOutputFuncBuilder(threshold), alertThreshold: threshold}
}

func (a *alertOutputProcessor) Transform(p []Payload) Payload {
	payload := convertToRequestVolumes(p)
	return a.transformFunc(payload)
}

func alertOutputFuncBuilder(threshold int) func(requestVolumes) Payload {
	return func(payload requestVolumes) Payload {
		numRequests := 0
		for _, val := range payload {
			numRequests += val.numRequests
		}
		endTime := payload[len(payload)-1].EndTime()
		duration := endTime - payload[0].StartTime()
		alertFiring := (float64(numRequests) / float64(duration)) > float64(threshold)
		return &volumeAlertStatus{alertFiring: alertFiring, time: endTime, volume: numRequests}
	}
}

func convertToLogItems(p []Payload) logItems {
	var payload logItems
	for _, val := range p {
		payload = append(payload, *(val.(*logItem)))
	}
	return payload
}

func convertToRequestVolumes(p []Payload) requestVolumes {
	var payload requestVolumes
	for _, val := range p {
		payload = append(payload, *(val.(*requestVolume)))
	}
	return payload
}

func getPayloadsTimestamp(p []Payload) timestamp {
	return timestamp{startTime: p[0].StartTime(), endTime: p[len(p)-1].EndTime()}
}
