package cmd

import (
	"strconv"
	"time"
)

type LogStat struct {
	writeFunc       func([][]string) HttpStats
	outFunc         func([][]string) int
	intervalSeconds int
}

func (l *LogStat) logStats(in <-chan logItem, write chan<- HttpStats, out chan<- requestVolume) {

	// read intervalSeconds worth of data from the channel
	// process it with outFunc and writeFunc then
	var minTimestamp int
	var data [][]string
	for x := range in {
		row := x.row
		data = append(data, row)
		if minTimestamp == 0 {
			minTimestamp, _ = strconv.Atoi(row[3])
		}
		thisTimestamp, _ := strconv.Atoi(row[3])
		if thisTimestamp < minTimestamp {
			minTimestamp = thisTimestamp
		}
		if thisTimestamp > minTimestamp+l.intervalSeconds {
			toWrite := l.writeFunc(data)
			write <- toWrite
			out <- requestVolume{numRequests: l.outFunc(data),
				err: nil, endTime: time.Unix(int64(thisTimestamp), 0)}
			minTimestamp = thisTimestamp
			data = nil
		}
	}
	close(write)
	close(out)

}
