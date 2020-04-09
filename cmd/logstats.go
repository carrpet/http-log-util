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
	IteratorKey() (int, error)
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
		//TODO: handle error
		it, _ := x.IteratorKey()

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
