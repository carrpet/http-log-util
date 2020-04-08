package cmd

import (
	"testing"
)

func TestLogStatsReadsIntervalWorthOfData(t *testing.T) {
	ls := LogStat{intervalSeconds: 10}
	logChan := make(chan logItem)
	go ls.logStats(logChan, nil, nil)

}
