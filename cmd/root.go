package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

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

	// get the flags for the high alert threshold, default 10 requests per second
	//alertThresholdPerSec := 10
	//alertFrequencySec := 120

	// open up log file for reading
	filereader, err := os.Open(args[0])
	if err != nil {
		return err
	}
	rowChan := make(chan logItem)
	requestsChan := make(chan requestVolume)
	statsChan := make(chan HttpStats)

	csvRows := newLogReader(filereader)
	doneChan := make(chan interface{})
	go csvRows.rows(rowChan)
	lf := LogFilter{interval: 10}
	statsTForm := statsTransform{tFunc: computeHTTPStats, out: statsChan}
	requestsTForm := requestVolumeTransform{tFunc: computeRequestVolume, out: requestsChan}
	go logItemsfilter(lf, rowChan, doneChan, statsTForm, requestsTForm)
	go func() {
		for x := range statsChan {
			fmt.Printf("Hits: %s for Section: %s\n", x.topHits[0].hits, x.topHits[0].section)
		}
	}()
	go func() {
		for y := range requestsChan {
			fmt.Printf("Num Requests: %d for End Time: %s\n", y.numRequests, time.Time.String(y.endTime))
		}
	}()

	//alertCfg := volumeAlertConfig{alertThreshold: alertThresholdPerSec,
	//	alertFrequency: alertFrequencySec}
	//alertCfg.requestVolumeAlert(requestVolChan, os.Stdout)
	<-doneChan
	return nil
}

// pipeline should contain a go channel that listens on a channel and
// consumes 10 seconds worth of data

// then another thing should listen on another channel from the first
// thing and consume 2 minutes worth of data
