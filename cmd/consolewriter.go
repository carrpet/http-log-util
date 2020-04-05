package cmd

import (
	"fmt"
	"time"
)

type ConsoleWriter struct {

}

type HttpStats struct {
	startTime time.Time
	endTime time.Time
	topHits []TopHitStat
}

type TopHitStat struct {
	section string
	hits string
}

func (hs *HttpStats) Print() {
	fmt.Printf("Top Hits For Time %s", hs.startTime.Format(time.RFC3339) )
	for _,th := range hs.topHits {
		fmt.Printf("{Section: %s, Number of Hits: %s",th.section, th.hits)
	}
	fmt.Println("/n")
}

func (cw *ConsoleWriter) Write(writeFrom <- chan HttpStats) {

	for x := range writeFrom {
		x.Print()
	}

} 