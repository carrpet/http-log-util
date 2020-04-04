package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

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
	_, err = csvReader.Read()

	//TODO: validate that the columns are in the expected order
	if err != nil {
		return err
	}
	buf := make([][]string, bufSize)
	for {
		read, err := readChunk(csvReader, buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		for i := 0; i < read; {
			data, _ := parseChunk(10, buf, read, i)

			//process the data here
			i = i + len(data)
		}

	}

	return nil
}

// returns the number of items read
func readChunk(from *csv.Reader, buf [][]string) (int, error) {

	i := 0
	rec, err := from.Read()
	if err != nil {
		return 0, err
	}
	buf[i] = rec
	minTimestamp, _ := strconv.Atoi(buf[i][3])
	maxTimestamp := minTimestamp

	// read in 2 minutes worth of data into buf
	for i = 1; maxTimestamp < minTimestamp+121 && i < len(buf); i++ {
		buf[i], err = from.Read()
		if err != nil {
			return i, err
		}
		tsStr, _ := strconv.Atoi(buf[i][3])
		if tsStr < minTimestamp {
			minTimestamp = tsStr
		} else {
			maxTimestamp = tsStr
		}
	}
	return i + 1, nil
}

// returns a slice that contains <seconds> seconds worth of data
func parseChunk(seconds int, buf [][]string, size int, start int) ([][]string, error) {

	minTimestamp, _ := strconv.Atoi(buf[start][3])
	thisTimestamp := minTimestamp
	i := start + 1
	for ; thisTimestamp < minTimestamp+11 && i < size; i++ {
		thisTimestamp, _ = strconv.Atoi(buf[i][3])
		if thisTimestamp < minTimestamp {
			minTimestamp = thisTimestamp
		}
	}
	return buf[start:i], nil
}

func computeTopHits(httpTraffic [][]string) []string {
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
	return []string{maxSection, strconv.Itoa(maxHits)}
}
