package cmd

import (
	"errors"
	"fmt"
	"os"

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

// Execute implements cobra's pattern for executing the command.
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

	//TODO: get user configured alert params
	alertThresholdPerSec := 10
	alertFrequencySec := 120

	// params to configure the log intervals
	logIntervalSec := 10

	// open up log file for reading
	filereader, err := os.Open(args[0])
	if err != nil {
		return err
	}

	// setup source and pass the params to it
	logSource := newCsvLogSource(filereader)
	// setup sink to be the log writers

	//setup and start the pipeline using source as source
	httpLogMonitor := NewPipeline(
		NewFanOutStage(
			[]Transformer{newRequestVolumeProcessor(), newHTTPStatsProcessor()},
			logIntervalSec),
		newStage(newAlertOutputProcessor(alertThresholdPerSec), alertFrequencySec))

	statsCh := make(chan Payload)
	sinkCh, _ := httpLogMonitor.Start(logSource, statsCh)

	// process the http stats
	go func() {
		for x := range statsCh {
			stat := x.(*httpStats)
			stat.Write(os.Stdout)
		}
	}()
	writeAlerts(sinkCh, os.Stdout)

	return nil
}
