package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "httpmonitorutil",
		Short: "HTTPMonitorUtil is a CLI utility to monitor http log files",
		Long: `A CLI utility to monitor http log files and gather metrics and useful 
				statistics about them.`,
		Args: argValidator,
		RunE: monitorCmd,
	}
	alertThresholdSec int
)

// Execute implements cobra's pattern for executing the command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVar(&alertThresholdSec, "alertThreshold", 10,
		"Defines the threshold in requests per seconds for which request volume alerts should fire")
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

// monitorCmd is the main routine to ingest an http log file
// and create and start a pipeline to write metrics and send alerts
func monitorCmd(cmd *cobra.Command, args []string) error {

	// hard coded params for the pipeline
	alertFrequencySec := 120
	logIntervalSec := 10

	// open up log file for reading
	filereader, err := os.Open(args[0])
	if err != nil {
		return err
	}

	// setup source and pass the params to it
	logSource := NewCSVLogSource(filereader)

	//setup and start the pipeline using csv file as source
	httpLogMonitor := NewPipeline(
		NewFanOutStage(
			[]Transformer{NewRequestVolumeProcessor(), NewHTTPStatsProcessor()},
			logIntervalSec),
		newStage(NewAlertProcessor(alertThresholdSec), alertFrequencySec))

	//channels that will be managed by this goroutine to retrieve output from the pipeline
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
