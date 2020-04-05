package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "httprequestutil",
		Short: "HttpRequestUtil is a CLI utility to monitor http log files",
		Long: `A CLI utility to monitor http log files and gather metrics and useful 
				statistics about them.`,
		Args: argValidator,
		RunE: monitorCmd,
	}
	logFile string
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

}

func argValidator(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires a file name argument")
	}
	info, err := os.Stat(args[0])
	if err != nil {
		return err
	}
	if info.IsDir() {
		return errors.New("Specified name should be a file not a directory")
	}
	return nil
}

// main routine to ingest log file and calculate and write
// metrics info
func monitorCmd(cmd *cobra.Command, args []string) error {

	// open up log file for reading
	filereader, err := os.Open(args[0])
	if err != nil {
		return err
	}
	csvRows := newLogReader(filereader)
	lStats := LogStat{
		writeFunc:       computeTopHits,
		intervalSeconds: 10,
	}
	statsChan := make(chan HttpStats)
	go lStats.logStats(csvRows.rows(), statsChan)
	for x := range statsChan {
		x.Print()
	}
	return nil
}

func computeTopHits(httpTraffic [][]string) HttpStats {
	hits := map[string]int{}
	for _, row := range httpTraffic {
		req := row[4]
		path := strings.Split(req, " ")
		section := "/" + strings.SplitN(path[1], "/", 3)[1]
		hits[section]++
	}

	//find the max hits section
	maxHits := 0
	var maxSection string
	for sect, h := range hits {
		if h > maxHits {
			maxHits = h
			maxSection = sect
		}
	}
	return HttpStats{topHits: []TopHitStat{{section: maxSection, hits: strconv.Itoa(maxHits)}}}
}

// pipeline should contain a go channel that listens on a channel and
// consumes 10 seconds worth of data

// then another thing should listen on another channel from the first
// thing and consume 2 minutes worth of data
