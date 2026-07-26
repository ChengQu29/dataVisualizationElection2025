package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type resultKey struct {
	districtNumber string
	districtName   string
	candidateName  string
	party          string
}

func main() {

	// make a map to hold the totals
	totals := make(map[resultKey]int64)

	fmt.Println("we are here0")
	// process one csv file
	err := processFile("../pollresults_resultatsbureauCanada/pollresults_resultatsbureau10001.csv", totals)
	if err != nil {
		panic(err)
	}

	// write the totals to a new csv file
	outputFile, err := os.Create("output_go.csv")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(bufio.NewWriter(outputFile))
	defer writer.Flush()

	// write the header row
	writer.Write([]string{"districtNumber", "districtName", "candidateName", "party", "totalVotes"})

	// write the totals to the csv file
	for key, total := range totals {
		writer.Write([]string{
			key.districtNumber,
			key.districtName,
			key.candidateName,
			key.party,
			strconv.FormatInt(total, 10),
		})
	}

}

func processFile(filePath string, totals map[resultKey]int64) error {
	// open the file
	f, err := os.Open(filePath)

	if err != nil {
		return err
	}
	defer f.Close()

	// create a new csv reader
	reader := csv.NewReader(bufio.NewReader(f))

	// read the header row
	header, err := reader.Read()

	//header is a string slice
	fmt.Println("Header type:", reflect.TypeOf(header))
	fmt.Println("Header:", header)

	if err != nil {
		return err
	}

	// create a map to hold the index of each column
	idx := make(map[string]int)
	for i, name := range header {
		idx[name] = i
	}

	// read the rest of the rows
	for {
		record, err := reader.Read()
		fmt.Println("Record:", record)
		if err != nil {
			break
		}

		// extract the relevant fields from the record
		districtNumber := record[idx["Electoral District Number/Numéro de circonscription"]]
		districtName := record[idx["Electoral District Name/Nom de la circonscription"]]
		candidateName := record[idx["Candidate’s Family Name/Nom de famille du candidat"]]
		party := record[idx["Political Affiliation Name_English/Appartenance politique_Anglais"]]
		votesStr := record[idx["Candidate Vote Count/Votes du candidat"]]

		// convert votes to int64
		votes, err := strconv.ParseInt(votesStr, 10, 64)
		if err != nil {
			return err
		}

		// create a key for the map
		key := resultKey{
			districtNumber: districtNumber,
			districtName:   districtName,
			candidateName:  candidateName,
			party:          party,
		}

		// add the votes to the total for this key
		totals[key] += votes
	}

	return nil
}
