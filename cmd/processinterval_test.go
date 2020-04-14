package cmd

import (
	"testing"
)

func TestProcessIntervalDoesNotReturnLastInterval(t *testing.T) {
	data := [6]mockPayload{{data: 0}, {data: 2}, {data: 2}, {data: 3}, {data: 5}, {data: 6}}
	expected := [6]bool{false, false, false, true, false, false}
	procFunc := processInterval(3)
	for i, val := range data {
		result := procFunc(val)
		if result != expected[i] {
			t.Errorf("ProcessInterval returned wrong value, expected %t received %t for i %d", expected[i], result, i)
		}
	}
}

func TestProcessIntervalTimestampsAreNotMonotonicIncreasing(t *testing.T) {
	data := [8]mockPayload{{data: 1549573868}, {data: 1549573867}, {data: 1549573868},
		{data: 1549573869}, {data: 1549573870}, {data: 1549573870}, {data: 1549573871}, {data: 1549573872}}
	expected := [8]bool{false, false, false, true, false, false, false, true}
	procFunc := processInterval(2)
	for i, val := range data {
		result := procFunc(val)
		if result != expected[i] {
			t.Errorf("ProcessInterval returned wrong value, expected %t received %t for i %d", expected[i], result, i)
		}
	}
}
