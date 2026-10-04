package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
)

func loadRecipient(filePath string, ch chan Recipient) error {

	defer close(ch)

	f, err := os.Open(filePath)

	if err != nil {
		return err
	}
	defer f.Close()
	
	r := csv.NewReader(f)
	r.TrimLeadingSpace = true

	records, err := r.ReadAll()

	if err != nil {
		return err
	}
	
	for _, record := range records[1:] { // 1: -> values after 1st index(0 idx)
		//fmt.Println(record)

		if len(record) < 2 {
    		log.Printf("skipping malformed row: %v", record)
    		continue
		}

		name := strings.TrimSpace(record[0])   // <- add back
		email := strings.TrimSpace(record[1])  // <- add back

		// sender -> consumer using unbuffered channel(queue)
		ch <- Recipient{
			Name: name,
			Email: email,
		}
	}

	return nil
}