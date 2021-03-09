package dataservice

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoURL string
var mongoDB = "geolocation-test"
var mongoCollection = "location"
var importFilePath = "/tmp/data_dump.csv"

var collection *mongo.Collection

func TestMain(m *testing.M) {
	os.Create(importFilePath)

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("mongo", "4.0.12", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	port := resource.GetPort("27017/tcp")
	mongoURL = fmt.Sprintf("mongodb://127.0.0.1:%s", port)

	if err := pool.Retry(func() error {
		var err error
		db, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURL))
		if err != nil {
			return err
		}
		collection = db.Database(mongoDB).Collection(mongoCollection)
		return db.Ping(context.Background(), nil)
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	collection.Drop(context.Background())
	os.Remove(importFilePath)

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestFetchLocation(t *testing.T) {
	ds, err := NewDataService(Settings{
		ImportFilePath:  importFilePath,
		MongoURL:        mongoURL,
		MongoDB:         mongoDB,
		MongoCollection: mongoCollection,
	})
	if err != nil {
		t.Errorf("Failed to initialize DataService: %v", err)
		return
	}

	_, err = ds.FetchLocation("8.8.8.8")
	if err == nil {
		t.Error("Fetched a location for non-existing IP record")
		return
	}

	country := "Netherlands"
	latitude := 50.0050
	longitude := 25.0025
	location := &Location{
		IPAddress:    "8.8.8.8",
		CountryCode:  "NL",
		Country:      country,
		City:         "Amsterdam",
		Latitude:     latitude,
		Longitude:    longitude,
		MysteryValue: 0,
	}

	collection.InsertOne(context.Background(), location)

	location, err = ds.FetchLocation("8.8.8.8")
	if err != nil {
		t.Errorf("Failed to fetch a location for existing IP record: %v", err)
		return
	}

	if location.Latitude != latitude || location.Longitude != longitude || location.Country != country {
		t.Error("Fetched a location for existing records but values do not match")
		return
	}
}

func TestImportData(t *testing.T) {
	testImportFileContent := `ip_address,country_code,country,city,latitude,longitude,mystery_value
200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346
144.116.254.249,,,,0,0,8050339844`

	err := ioutil.WriteFile(importFilePath, []byte(testImportFileContent), 0644)
	if err != nil {
		t.Errorf("Failed to write to the import file: %v", err)
		return
	}

	ds, err := NewDataService(Settings{
		ImportFilePath:  importFilePath,
		MongoURL:        mongoURL,
		MongoDB:         mongoDB,
		MongoCollection: mongoCollection,
	})
	if err != nil {
		t.Errorf("Failed to initialize DataService: %v", err)
		return
	}

	statistics, err := ds.ImportData(false)
	if err != nil {
		t.Errorf("Failed to import data: %v", err)
		return
	}

	if statistics.AcceptedRecordCount != 1 || statistics.RejectedRecordCount != 1 {
		t.Error("ImportData statistics are not as expected")
		return
	}

	fetchedResult := collection.FindOne(context.Background(), bson.D{primitive.E{Key: "ipAddress", Value: "200.106.141.15"}})
	if fetchedResult.Err() != nil {
		t.Errorf("Failed to fetch the result of data import: %v", fetchedResult.Err())
		return
	}

	fetchedResult = collection.FindOne(context.Background(), bson.D{primitive.E{Key: "ipAddress", Value: "144.116.254.249"}})
	if fetchedResult.Err() == nil {
		t.Errorf("Successfully fetched a record that shouldn't have been persisted")
		return
	}
}
