package dataservice

import (
	"context"
	"errors"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Settings holds the values that are used for initializing an instance of the library
type Settings struct {
	ImportFilePath  string
	MongoURL        string
	MongoDB         string
	MongoCollection string
}

// DataService is a handle where library configuration is stored and library methods can be accessed
type dataService struct {
	settings        Settings
	importFile      *os.File
	mongoClient     *mongo.Client
	mongoCollection *mongo.Collection
}

// NewDataService initializes a new instance of the library and returns a handle
func NewDataService(settings Settings) (*dataService, error) {
	if settings.ImportFilePath == "" || settings.MongoURL == "" || settings.MongoDB == "" || settings.MongoCollection == "" {
		return nil, errors.New("Not all setting parameters are provided")
	}

	file, err := os.Open(settings.ImportFilePath)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(settings.MongoURL))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	collection := client.Database(settings.MongoDB).Collection(settings.MongoCollection)

	err = createUniqueIndexOnIPAddress(collection)
	if err != nil {
		return nil, err
	}

	return &dataService{settings, file, client, collection}, nil
}
