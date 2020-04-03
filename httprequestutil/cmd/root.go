package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

const (
	bufSize = 4096
)

func most_hits(rec [][]string) (*io.Writer, error) {
	return nil, nil
}

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
	//os.

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
	filereader, err := os.Open(logFile)
	if err != nil {
		return err
	}
	csvReader := csv.NewReader(filereader)
	rec, err := csvReader.Read()
	if err != nil {
		panic(err)
	}
	buf := make([][]string, bufSize)
	i := 0
	for rec != nil {
		fmt.Printf("Reading record: %s\n", rec)
		rec, err = csvReader.Read()
		if err != nil {
			break
			//panic(err)
		}
		minTimestamp, _ := strconv.Atoi(rec[3])
		maxTimestamp := minTimestamp
		// read in 2 minutes worth of data into buf
		for maxTimestamp < minTimestamp+121 {
			thisData, err := csvReader.Read()
			if err != nil {
				panic(err)
			}
			buf[i] = thisData
			tsStr, _ := strconv.Atoi(buf[i][3])
			if tsStr < minTimestamp {
				minTimestamp = tsStr
			} else {
				maxTimestamp = tsStr
			}
			i++
		}

		// the  data in 10 second chunks
		startPtr := 0
		endPtr := 0

		for endPtr < i {
			chunkTimestamp := minTimestamp
			for chunkTimestamp < minTimestamp+11 {
				thisStr, _ := strconv.Atoi(buf[endPtr][3])
				if thisStr < minTimestamp {
					minTimestamp = thisStr
				} else {
					chunkTimestamp = thisStr
				}
				endPtr++
			}

			// now process the slice of data
			most_hits(buf[startPtr:endPtr])

			//write the data

			startPtr = endPtr
			minTimestamp, err = strconv.Atoi(buf[endPtr][3])

		}
	}
	return nil
}
