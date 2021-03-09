package dataservice

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// createUniqueIndexOnIPAddress creates a MongoDB unique index on `ipAddress` field so that duplicate entries can be rejected
func createUniqueIndexOnIPAddress(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{Keys: bson.M{"ipAddress": 1}, Options: options.Index().SetUnique(true)}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}
