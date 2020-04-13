package cmd

//FlowRegulator does something
type FlowRegulator interface {
	ProcessInterval() func(...int) bool
}

// LogFilter will implement FlowRegulator
type LogFilter struct {
	interval int //input control
}

// ProcessInterval returns a function that when called
// takes a sequence of non-monotonically increasing integers
// and returns true if the number of arguments called equals or
// exceeds the configured interval.
func (lf LogFilter) ProcessInterval() func(...int) bool {
	minTime := -1
	return func(ints ...int) bool {
		for _, val := range ints {
			if minTime == -1 {
				minTime = val
			}
			if val < minTime {
				minTime = val
			}
			if val >= minTime+lf.interval {
				minTime = -1
				return true
			}
		}
		return false
	}
}
