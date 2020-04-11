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
	Write()
}
type Iterable interface {
	IteratorKey() (int, error)
}

func (li *logItem) IteratorKey() (int, error) {
	ts, err := strconv.Atoi(li.row[date])
	if err != nil {
		// TODO
		fmt.Errorf("Error retrieving iterator key for logItem: %s", err.Error())
	}
	return ts, nil

}

//FlowRegulator does something
type FlowRegulator interface {
	ProcessInterval() func(...int) bool
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

// ProcessInterval returns a function that when called
// takes a sequence of non-monotonically increasing integers
// and returns true if the number of arguments called equals or
// exceeds the configured interval.
func (lf LogFilter) ProcessInterval() func(...int) bool {
	minTime := 0
	return func(ints ...int) bool {
		for _, val := range ints {
			if minTime == 0 {
				minTime = val
			}
			if val < minTime {
				minTime = val
			}
			if val > minTime+lf.interval {
				minTime = 0
				return true
			}
		}
		return false
	}
}

type LogItems []logItem

/*
type Transformable interface {
	Transform() Iterable
	TransformForWrite() Writable
}
*/
/*
type statsTransform struct {
	tFunc func([]Iterable) HttpStats
	out   chan<- HttpStats
}

func (s statsTransform) SendTransformation(items []Iterable) {
	s.out <- s.tFunc(items)
}
*/

type requestVolumeTransform struct {
	tFunc func(LogItems) requestVolume
	out   chan<- requestVolume
}

type alerts struct {
}

type alertTransform struct {
	tFunc func([]requestVolume) alerts
	out   chan<- alerts
}

type stageConfig struct {
	transformer    Transformer
	writeTransform Transformer
}

func (s stageConfig) logItemsStage(fr FlowRegulator, done chan<- interface{}, p StageParams) {

	processFunc := fr.ProcessInterval()
	var buf []Payload
	for x := range p.Input() {
		buf = append(buf, x)
		//TODO: handle error
		it, _ := x.IteratorKey()

		iVal := processFunc(it)
		if !iVal {
			continue
		}
		result, err := s.transformer.Transform(buf)
		if err != nil {
			fmt.Errorf("Do something %s", err)
		}
		p.Output() <- result
		buf = nil
	}
	done <- struct{}{}
}

type Transformable interface {
	Transform() Transformable
}

type LIS struct {
}

func (l LIS) Transform() Transformable {

}

// takes a logItem, flowRegulator and outputs a Transformable collection
func testFilter(fr FlowRegulator, in <-chan logItem, out chan<- Transformable) {

	var items LogItems
	for x := range in {
		items = append(items, x)
	}
	out <- items.Transform()

}

type testTForm struct {
	tForm func(LogItems) requestVolume
	out   chan requestVolume
}

func (t testTForm) Transform(li LogItems) Transformable {
	return t.tForm(li)
}

type httpStatsProcessor struct {
}
