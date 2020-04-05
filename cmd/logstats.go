package cmd

import "strconv"

type LogStat struct {
	writeFunc       func([][]string) HttpStats
	outFunc         func([][]string) int
	intervalSeconds int
}

func (l *LogStat) logStats(in <-chan []string, write chan<- HttpStats) {

	// read intervalSeconds worth of data from the channel
	// process it with outFunc and writeFunc then
	var minTimestamp int
	var data [][]string
	for x := range in {
		data = append(data, x)
		if minTimestamp == 0 {
			minTimestamp, _ = strconv.Atoi(x[3])
		}
		thisTimestamp, _ := strconv.Atoi(x[3])
		if thisTimestamp < minTimestamp {
			minTimestamp = thisTimestamp
		}
		if thisTimestamp > minTimestamp+l.intervalSeconds {
			toWrite := l.writeFunc(data)
			write <- toWrite
			//out <- l.outFunc(data)
			minTimestamp = thisTimestamp
			data = nil
		}
	}
	close(write)
	//close(out)

}
