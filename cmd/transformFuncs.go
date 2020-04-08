package cmd

import (
	"strconv"
	"strings"
	"time"
)

func computeHTTPStats(items LogItems) HttpStats {
	hits := map[string]int{}
	for _, val := range items {
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

func computeRequestVolume(items LogItems) requestVolume {
	endTime, err := strconv.Atoi(items[len(items)-1].row[date])
	if err != nil {
		return requestVolume{err: err}
	}
	return requestVolume{numRequests: len(items),
		err: nil, endTime: time.Unix(int64(endTime), 0)}
}
