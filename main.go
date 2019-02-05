package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/smartystreets/scanners/csv"
)

func main() {
	start := time.Now()
	csvIn, err := os.Open("test_data_ridiculous.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvIn.Close()

	csvOut, _ := os.Create("myout.csv")
	writer := csv.NewWriter(csvOut, ';')

	ch := scanCSV(csvIn)
	err = writer.WriteStream(ch)
	if err != nil {
		log.Panic(err)
	}

	elapsed := time.Since(start)
	log.Println("Time elapsed: ", elapsed)

	os.Remove(csvOut.Name())
}

func scanCSV(rc io.Reader) (ch chan []string) {
	// could pass delimiter here
	// allow to define a header row
	ch = make(chan []string, 10)
	go func() {
		scanner := csv.NewScanner(rc, csv.Comma(','), csv.Comment('#'), csv.TrimLeadingSpace(true))
		defer close(ch)

		for scanner.Scan() {
			if err := scanner.Error(); err != nil {
				log.Panic(err)
			} else {
				// take row and build csv string
				// scanner.Record() is a slice of strings
				// log.Println("Record: ", scanner.Record())
				ch <- scanner.Record()
			}
		}
	}()
	return
}
