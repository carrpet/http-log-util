package cmd

import (
	"fmt"
	"strconv"
)

// represents the fields of the log file
const (
	rh       = iota
	rfc931   = iota
	authuser = iota
	date     = iota
	req      = iota
	status   = iota
	numBytes = iota
)

type Errorable interface {
	Error() bool
}
type Writable interface {
	WriteLog()
}
type Transformable interface {
	//Transform(func(Transformable) Errorable)
	IteratorKey() (int, error)
	//Transforms() func() Errorable
}

func (li logItem) IteratorKey() (int, error) {
	ts, err := strconv.Atoi(li.row[date])
	if err != nil {
		// TODO
		fmt.Errorf("Error retrieving iterator key for logItem: %s", err.Error())
	}
	return ts, nil

}

type FlowRegulator interface {
	ProcessInterval(int) func(...int) *Interval
}

// LogFilter will implement FlowRegulator
type LogFilter struct {
	interval int //input control
}

type Interval struct {
	start int
	end   int
}
type FilterConfig struct {
	intervalSeconds int //how much data should be aggregated at a time
}

// expects a non-monotonically increasing
// sequence and returns a struct with containing
// the start and end

func (lf LogFilter) ProcessInterval(start int) func(...int) *Interval {
	minTime := 0
	return func(ints ...int) *Interval {
		for _, val := range ints {
			if minTime == 0 {
				minTime = val
			}
			if val < minTime {
				minTime = val
			}
			if val > minTime+lf.interval {
				result := &Interval{start: minTime, end: val}
				minTime = 0
				//data = nil
				return result
			}
		}
		return nil
	}
}

type LogStat struct {
	data            [][]string
	writeFunc       func([][]string) HttpStats
	outFunc         func([][]string) int
	intervalSeconds int
}

//func (l *LogStat) lS(in <-chan logItem, write ...chan<- HttpStats)

/*
func (l *LogStat) logStats(in <-chan logItem, write chan<- HttpStats, out chan<- requestVolume) {

	// read intervalSeconds worth of data from the channel
	// process it with outFunc and writeFunc
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
*/

/*
func (l *LogStat) Transform(tFunc func ([][]string) Errorable, data [][]string) func(chan<-Errorable) {
	// define a function to apply to transform
	toTransform := data
	return func (out chan<-Errorable) {
		transformed := tFunc(toTransform)
		out <- transformed
	}
}
*/

type LogItems []logItem

type Sendable interface {
	SendTransformation(LogItems)
}
type statsTransform struct {
	tFunc func(LogItems) HttpStats
	out   chan<- HttpStats
}

func (s statsTransform) SendTransformation(items LogItems) {
	s.out <- s.tFunc(items)
}

type requestVolumeTransform struct {
	tFunc func(LogItems) requestVolume
	out   chan<- requestVolume
}

func (r requestVolumeTransform) SendTransformation(items LogItems) {
	r.out <- r.tFunc(items)
}

func logItemsfilter(fr FlowRegulator, in <-chan logItem, done chan<- interface{}, out ...Sendable) {

	processFunc := fr.ProcessInterval(0)
	var buf LogItems
	for x := range in {
		buf = append(buf, x)
		it, _ := x.IteratorKey()
		//TODO: handle iterator
		/* this could be abstracted into getSequenceItem
		row := x.row
		data = append(data, row)
		ts, err := strconv.Atoi(row[date])
		if err != nil {
			// TODO
			fmt.Println("Error parsing data")
		}
		*/
		iVal := processFunc(it)
		if iVal == nil {
			continue
		}
		for _, val := range out {
			val.SendTransformation(buf)
		}
		buf = nil
	}
	done <- struct{}{}
}
