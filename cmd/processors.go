package cmd

import (
	"fmt"
	"strconv"
	"time"
)

type LogItems []logItem

type logItem struct {
	row []string
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
	var payload LogItems
	for _, val := range p {
		payload = append(payload, *(val.(*logItem)))
	}
	endTime, err := strconv.Atoi(payload[len(payload)-1].row[date])
	if err != nil {
		return &requestVolume{err: err}
	}
	return &requestVolume{numRequests: len(payload),
		err: nil, endTime: time.Unix(int64(endTime), 0)}
}
