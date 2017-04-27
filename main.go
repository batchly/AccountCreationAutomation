package main

import (
	"accountcreationautomation/channels"
	"accountcreationautomation/logger"
	"accountcreationautomation/utils"
	"encoding/csv"
	"log"
	"os"
)

func main() {
	//Getting the API Access key and Seretkey from Commandline
	logger.Set()
	logger.Get().Infoln("Application started")
	defer logger.Get().Infoln("Application ends")

	ApiAccessKey := os.Args[1]
	ApiSecretKey := os.Args[2]
	DomainURL := os.Args[3]
	Filename := os.Args[4]

	records, err := utils.ReadCSV(Filename)
	if err != nil {
		log.Fatal(err)
	}

	// write results to a new csv
	outfile, err := os.Create("results.csv")
	if err != nil {
		log.Fatal("Unable to access results.csv,Please close the file and run the program")
	}
	defer outfile.Close()

	writer := csv.NewWriter(outfile)
	defer writer.Flush()
	//Getting the data from the crendentials.csv

	channels.New(records, ApiAccessKey, ApiSecretKey, DomainURL, writer)

}
