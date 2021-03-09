package dataservice

import (
	"testing"
)

func TestParseRecords(t *testing.T) {
	records := [][]string{
		{"8.8.8.8", "NL", "Netherlands", "Amsterdam", "-50.050", "25.025", "123"},
		{"9.9.9.9", "AZ", "Azerbaijan", "Baku", "Not a float", "25.025", "123"},
		{"11.11.11.11", "US", "123 Not a country", "123 Not a city", "50.050", "25.025", "123"},
	}

	locations, statistics := parseRecords(records)

	if len(locations) != 1 {
		t.Error("Parsed record count is more than expected")
		return
	}

	if statistics.AcceptedRecordCount != 1 || statistics.RejectedRecordCount != 2 {
		t.Error("Returned parse statistics does not contain the expected values")
		return
	}

	if locations[0].IPAddress != "8.8.8.8" || locations[0].Country != "Netherlands" || locations[0].Latitude != -50.050 {
		t.Error("Parsed record does not contain the expected values")
		return
	}
}
