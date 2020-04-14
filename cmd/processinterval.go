package cmd

// ProcessInterval returns a function that when called
// takes a sequence of non-monotonically increasing integers
// and returns true if the number of arguments called equals or
// exceeds the configured interval.
func processInterval(interval int) func(p Payload) bool {
	minTime := -1
	elapsed := 0
	return func(p Payload) bool {
		if minTime == -1 || p.StartTime() < minTime {
			minTime = p.StartTime()
		}
		elapsed = p.EndTime() - minTime
		if elapsed >= minTime+interval {
			minTime = -1
			elapsed = 0
			return true
		}
		return false
	}
}
