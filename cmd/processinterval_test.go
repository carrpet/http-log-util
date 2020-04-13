package cmd

import (
	"testing"
)

func TestProcessIntervalDoesNotReturnLastInterval(t *testing.T) {
	lf := LogFilter{interval: 3}
	data := [6]int{0, 2, 2, 3, 5, 6}
	expected := [6]bool{false, false, false, true, false, false}
	procFunc := lf.ProcessInterval()
	for i, val := range data {
		result := procFunc(val)
		if result != expected[i] {
			t.Errorf("ProcessInterval returned wrong value, expected %t received %t for i %d", expected[i], result, i)
		}
	}
}

func TestProcessIntervalTimestampsAreNotMonotonicIncreasing(t *testing.T) {
	lf := LogFilter{interval: 2}
	data := [8]int{1549573868, 1549573867, 1549573868, 1549573869, 1549573870, 1549573870, 1549573871, 1549573872}
	expected := [8]bool{false, false, false, true, false, false, false, true}
	procFunc := lf.ProcessInterval()
	for i, val := range data {
		result := procFunc(val)
		if result != expected[i] {
			t.Errorf("ProcessInterval returned wrong value, expected %t received %t for i %d", expected[i], result, i)
		}
	}
} 