package cmd

import (
	"net/http"
	"strconv"
	"strings"
)

// logItems represent objects coming from the csv source
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

// httpStats represents stats that are
// collected about the http log data
type httpStats struct {
	ts               timestamp
	topHits          []topHitStat
	unsuccessfulReqs []nonSuccessHits
}

type nonSuccessHits struct {
	section string
	hits    int
}

func (s *httpStats) StartTime() int { return s.ts.startTime }
func (s *httpStats) EndTime() int   { return s.ts.endTime }

type topHitStat struct {
	section string
	hits    string
}

type sectionStats struct {
	hits               int
	nonSuccessRequests int
}

// HTTPStatsProcessor represents configuration for
// a processor to implement gathering http log statistics for a pipeline stage.
type HTTPStatsProcessor struct {
	transformFunc func(logItems, timestamp) Payload
}

// NewHTTPStatsProcessor returns an httpStatsProcessor
// that can be used in processing a pipeline stage.
func NewHTTPStatsProcessor() *HTTPStatsProcessor {
	return &HTTPStatsProcessor{transformFunc: httpStatsTransformFunc}
}

// Transform applies the httpStatsProcessor
// transform func to convert LogItems to requestVolumes
// that are to be sent downstream.
func (hs *HTTPStatsProcessor) Transform(p []Payload) Payload {
	ts := getPayloadsTimestamp(p)
	payload := convertToLogItems(p)
	return hs.transformFunc(payload, ts)
}

// httpStatsTransformFunc computes meaningful statistics about
// a collection of logItems to obtain http metrics and alert data
// for downstream stages.
func httpStatsTransformFunc(p logItems, ts timestamp) Payload {

	stats := map[string]*sectionStats{}
	for _, val := range p {
		status, _ := strconv.Atoi(val.row[status])
		req := val.row[req]
		path := strings.Split(req, " ")
		section := "/" + strings.SplitN(path[1], "/", 3)[1]
		if _, ok := stats[section]; !ok {
			stats[section] = &sectionStats{hits: 1}
		} else {
			stats[section].hits++
		}
		if status != http.StatusOK {
			stats[section].nonSuccessRequests++
		}
	}

	//find the max hits section
	maxHits := 0
	var maxSection string
	badHits := []nonSuccessHits{}
	for section, s := range stats {
		badHits = append(badHits, nonSuccessHits{section: section, hits: s.nonSuccessRequests})
		if s.hits > maxHits {
			maxHits = s.hits
			maxSection = section
		}
	}
	return &httpStats{ts: ts, unsuccessfulReqs: badHits,
		topHits: []topHitStat{{section: maxSection, hits: strconv.Itoa(maxHits)}}}

}

// requestVolume represents
// the number of http requests recorded
// for a given time interval
type requestVolume struct {
	numRequests int
	ts          timestamp
}

type requestVolumes []requestVolume

func (r *requestVolume) StartTime() int { return r.ts.startTime }
func (r *requestVolume) EndTime() int   { return r.ts.endTime }

// RequestVolumeProcessor represents configuration for
// a processor to implement counting number of http requests for a pipeline stage.
type RequestVolumeProcessor struct {
	transformFunc func(logItems, timestamp) Payload
}

// NewRequestVolumeProcessor returns a requestVolumeProcessor
// that can be used in processing a pipeline stage.
func NewRequestVolumeProcessor() *RequestVolumeProcessor {
	return &RequestVolumeProcessor{transformFunc: requestVolumeTransformFunc}
}

// Transform applies the requestProcessor
// transform func to convert logItems to requestVolumes
// that are to be sent downstream.
func (r *RequestVolumeProcessor) Transform(p []Payload) Payload {
	ts := getPayloadsTimestamp(p)
	payload := convertToLogItems(p)
	return r.transformFunc(payload, ts)
}

func requestVolumeTransformFunc(payload logItems, ts timestamp) Payload {
	return &requestVolume{numRequests: len(payload), ts: ts}
}

// volumeAlertStatus represents the alert state for a time interval
type volumeAlertStatus struct {
	alertFiring bool
	time        int
	volume      int
}

func (v *volumeAlertStatus) StartTime() int { return v.time }
func (v *volumeAlertStatus) EndTime() int   { return v.time }

// AlertProcessor represents configuration for
// a processor to implement http request volume alerts for a pipeline stage.
type AlertProcessor struct {
	alertThreshold int
	transformFunc  func(requestVolumes) Payload
}

// NewAlertProcessor returns an alertProcessor configured with the given
// alert threshold that can be used in processing a pipeline stage.
func NewAlertProcessor(threshold int) *AlertProcessor {
	return &AlertProcessor{transformFunc: alertFuncBuilder(threshold), alertThreshold: threshold}
}

// Transform applies the alertProcessor
// transform func to convert requestVolumes to voluemAlertStatus
// that are to be sent downstream.
func (a *AlertProcessor) Transform(p []Payload) Payload {
	payload := convertToRequestVolumes(p)
	return a.transformFunc(payload)
}

// alertFuncBuilder returns a function to compute when volume alerts
// should fire based on the configured threshold.
func alertFuncBuilder(threshold int) func(requestVolumes) Payload {
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
