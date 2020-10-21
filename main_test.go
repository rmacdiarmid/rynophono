package main

//i am ryan branch

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

// AddNumbersToFile opens filename, reads each line, gets name and address, retrieves phone from places API and appends to outfile
func AddNumbersToFile(filename string) error {
	inf, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer inf.Close()
	incsv := csv.NewReader(inf)

	outf, err := os.OpenFile(filename+"_out.csv", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer outf.Close()
	outcsv := csv.NewWriter(outf)
	defer outcsv.Flush()

	b := Borrower{}
	currLine := 0
	for {
		record, err := incsv.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		currLine++
		if currLine == 1 { //headers
			outcsv.Write(AddPlaceToRecord(record, "PlaceID", "Phone", "Website"))
			continue
		}
		b.Name = record[2]
		b.Street = record[3]
		b.City = record[4]
		b.State = record[5]
		b.Zip = record[6]

		// get phone from API
		pr, err := GetPlacesId(&b)
		if err != nil {
			return err
		}
		if pr.Status != "OK" {
			outcsv.Write(AddPlaceToRecord(record, "", "", ""))
			continue
		}
		fmt.Println("getting place", pr.Candidates[0].PlaceID)

		pd, err := GetPlacesDetail(pr)
		if err != nil {
			return err
		}
		// append results to outfile
		outcsv.Write(AddPlaceToRecord(record, pr.Candidates[0].PlaceID, pd.Result.FormattedPhoneNumber, pd.Result.Website))
		outcsv.Write(record)
		fmt.Println(b)
	}
	return nil
}

func TestTextQuery(t *testing.T) {
	b := Borrower{Street: "123 Main", City: "Loomis", State: "CA", Zip: "95650"}
	if b.ToTextQuery() != "123 Main, Loomis CA 95650" {
		log.Fatal("ToTextQuery != ", b.ToTextQuery(), "123 Main, Loomis CA 95650")
	}
}

func TestCSV(t *testing.T) {
	log.Println("api key", gapikey)
	err := AddNumbersToFile("testdata/numbers.csv")
	if err != nil {
		log.Fatal("We died mom", err)
	}
	log.Println("Mom, we made it...")

}
