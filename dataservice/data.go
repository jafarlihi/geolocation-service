package dataservice

import (
	"encoding/csv"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"
)

// readCSVData reads the CSV file and returns the records in an array of string arrays
func readCSVData(file *os.File) ([][]string, error) {
	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	records, err := reader.ReadAll()

	if err != nil {
		return nil, err
	}

	file.Seek(0, 0)

	return records, nil
}

// isCountryCode validates country code values by checking if they are 2 capital letters
var isCountryCode = regexp.MustCompile(`^[A-Z]{2}$`).MatchString

// isCountryOrCity validates country and city values by checking if they are non-numeric strings with at least length of 1
var isCountryOrCity = regexp.MustCompile(`^[\D]+$`).MatchString

// parseRecords takes in CSV data as array of string arrays and decodes them into Location objects
// Values are checked for sanity and accepted/rejected count statistics are returned
func parseRecords(records [][]string) ([]Location, ImportStatistics) {
	var locations []Location
	acceptedRecordCount := 0
	rejectedRecordCount := 0

	for _, record := range records {
		if net.ParseIP(record[0]) == nil {
			rejectedRecordCount += 1
			continue
		}
		if !isCountryCode(record[1]) {
			rejectedRecordCount += 1
			continue
		}
		if !isCountryOrCity(record[2]) || !isCountryOrCity(record[3]) {
			rejectedRecordCount += 1
			continue
		}
		latitude, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			rejectedRecordCount += 1
			continue
		}
		longitude, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			rejectedRecordCount += 1
			continue
		}
		mysteryValue, err := strconv.ParseInt(record[6], 0, 64)
		if err != nil {
			mysteryValue = 0
		}

		location := Location{
			IPAddress:    record[0],
			CountryCode:  record[1],
			Country:      record[2],
			City:         record[3],
			Latitude:     latitude,
			Longitude:    longitude,
			MysteryValue: mysteryValue,
		}

		acceptedRecordCount += 1
		locations = append(locations, location)
	}

	return locations, ImportStatistics{time.Duration(0), acceptedRecordCount, rejectedRecordCount}
}
