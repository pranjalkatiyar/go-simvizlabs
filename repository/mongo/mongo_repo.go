package mongoRepo

import (
	"context"
	"simvizlab-backend/infra/database"
	"simvizlab-backend/infra/logger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultDatabaseName = "userdb"
	defaultTimeout      = 10 * time.Second
)

// Save inserts a new document into MongoDB
func Save(collection string, model interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	coll := database.MongoClient.Database(defaultDatabaseName).Collection(collection)
	_, err := coll.InsertOne(ctx, model)
	if err != nil {
		logger.Errorf("error saving data to MongoDB: %v", err)
		return err
	}
	return nil
}

// Get retrieves all documents from a collection
func Get(collection string, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	coll := database.MongoClient.Database(defaultDatabaseName).Collection(collection)
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, result)
}

// GetOne retrieves the last document from a collection
func GetOne(collection string, filter bson.M, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	coll := database.MongoClient.Database(defaultDatabaseName).Collection(collection)
	opts := options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}})

	err := coll.FindOne(ctx, filter, opts).Decode(result)
	if err != nil {
		return err
	}
	return nil
}

// Update updates a document in MongoDB
func Update(collection string, filter bson.M, update interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	coll := database.MongoClient.Database(defaultDatabaseName).Collection(collection)
	_, err := coll.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a document from MongoDB
func Delete(collection string, filter bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	coll := database.MongoClient.Database(defaultDatabaseName).Collection(collection)
	_, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
