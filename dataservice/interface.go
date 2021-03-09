package dataservice

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IDataService is an interface for dataService
type IDataService interface {
	FetchLocation(ipAddress string) (*Location, error)
	ImportData(dropExistingRecords bool) (*ImportStatistics, error)
}

// Location is the main struct used for exposing stored data to consumers
type Location struct {
	IPAddress    string  `bson:"ipAddress"`
	CountryCode  string  `bson:"countryCode"`
	Country      string  `bson:"country"`
	City         string  `bson:"city"`
	Latitude     float64 `bson:"latitude"`
	Longitude    float64 `bson:"longitude"`
	MysteryValue int64   `bson:"mysteryValue"`
}

// ImportStatistics is used for holding the information about a completed data import process
type ImportStatistics struct {
	Duration            time.Duration
	AcceptedRecordCount int
	RejectedRecordCount int
}

// FetchLocation is one of two main interface methods to the library, it takes in an IP address and attempts to return the stored location information about it
func (ds dataService) FetchLocation(ipAddress string) (*Location, error) {
	var location Location

	err := ds.mongoCollection.FindOne(context.Background(), bson.D{{"ipAddress", ipAddress}}).Decode(&location)
	if err != nil {
		return nil, err
	}

	return &location, nil
}

// ImportData is one of the two main interface methods to the library, it reads the configured CSV file and imports it to the database
func (ds dataService) ImportData(dropExistingRecords bool) (*ImportStatistics, error) {
	startTime := time.Now()

	records, err := readCSVData(ds.importFile)
	if err != nil {
		return nil, err
	}

	if dropExistingRecords {
		ds.mongoCollection.Drop(context.Background())
		err = createUniqueIndexOnIPAddress(ds.mongoCollection)
		if err != nil {
			return nil, err
		}
	}

	locations, importStatistics := parseRecords(records)

	var locationsInterface []interface{}
	for i, location := range locations {
		locationsInterface = append(locationsInterface, location)
		// This is kinda hacky but we run InsertMany after every 50000th record because there seems to be a bug in MongoDB
		// driver where it nevers returns from InsertMany when too many records are passed in and `ordered` option is set
		// to `false` -- not even after `db.currentOp()` stops reporting any running operations.
		if i%50000 == 0 {
			ds.mongoCollection.InsertMany(context.Background(), locationsInterface, options.InsertMany().SetOrdered(false))
			locationsInterface = locationsInterface[:0]
		}
	}
	ds.mongoCollection.InsertMany(context.Background(), locationsInterface, options.InsertMany().SetOrdered(false))

	elapsedTime := time.Since(startTime)
	importStatistics.Duration = elapsedTime

	return &importStatistics, nil
}
