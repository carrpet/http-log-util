package cmd

import (
	"fmt"
	"time"
)

type ConsoleWriter struct {
}

type HttpStats struct {
	startTime time.Time
	endTime   time.Time
	topHits   []TopHitStat
}

//IteratorKey currently does nothing
func (hs HttpStats) IteratorKey() (int, error) {
	return 0, nil
}

type TopHitStat struct {
	section string
	hits    string
}

// Write writes the stats to the console.
func (hs HttpStats) Write() {
	fmt.Printf("Top Hits For Time %s", hs.startTime.Format(time.RFC3339))
	for _, th := range hs.topHits {
		fmt.Printf("{Section: %s, Number of Hits: %s", th.section, th.hits)
	}
	fmt.Println("/n")
}

func (cw *ConsoleWriter) Write(writeFrom <-chan HttpStats) {

	for x := range writeFrom {
		x.Print()
	}

}
