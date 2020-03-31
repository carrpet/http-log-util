package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {

	log := os.Args[1]
	// open up log file for reading
	filereader, err := os.Open(log)
	if err != nil {
		panic("Couldn't open the file!!")
	}
	csvReader := csv.NewReader(filereader)
	rec, err := csvReader.Read()
	if err != nil {
		panic(err)
	}
	for rec != nil {
		fmt.Printf("Reading record: %s\n", rec)
		rec, err = csvReader.Read()
		if err != nil {
			panic(err)
		}
	}
}
