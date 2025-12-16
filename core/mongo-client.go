package core

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoClient struct {
	client *mongo.Client
}

func NewMongoClient() *MongoClient {
	uri := buildMongoURI()
	// Uses the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// Defines the options for the MongoDB client
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Creates a new client and connects to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}
	// Sends a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.M{"ping": 1}).Decode(&result); err != nil {
		panic(err)
	}
	return &MongoClient{
		client: client,
	}
}

func validateEnvVariables(vars []string) {
	for _, v := range vars {
		if os.Getenv(v) == "" {
			log.Fatalf("You must set your '%s' environment variable.", v)
		}
	}
}

func buildMongoURI() string {
	validateEnvVariables([]string{
		"MONGODB_HOST",
		"MONGODB_PORT",
		"MONGODB_DB",
		"MONGODB_USER",
		"MONGODB_PASSWORD",
		"MONGODB_AUTH_SOURCE",
	})
	host := os.Getenv("MONGODB_HOST")
	port := os.Getenv("MONGODB_PORT")
	database := os.Getenv("MONGODB_DB")
	user := os.Getenv("MONGODB_USER")
	password := os.Getenv("MONGODB_PASSWORD")
	authSource := os.Getenv("MONGODB_AUTH_SOURCE")

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s", user, password, host, port, database, authSource)
	return uri
}

// Helper method to get a collection by name
func (mc *MongoClient) getCollection(collectionName string) *mongo.Collection {
	databaseName := os.Getenv("MONGODB_DB")
	return mc.client.Database(databaseName).Collection(collectionName)
}

// InsertOne inserts a single document into the specified collection
func (mc *MongoClient) InsertOne(collectionName string, document interface{}) (*mongo.InsertOneResult, error) {
	collection := mc.getCollection(collectionName)
	result, err := collection.InsertOne(context.TODO(), document)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %v", err)
	}
	return result, nil
}

// FindOne finds a single document in the specified collection
func (mc *MongoClient) FindOne(collectionName string, filter interface{}) *mongo.SingleResult {
	collection := mc.getCollection(collectionName)
	result := collection.FindOne(context.TODO(), filter)
	return result
}

// FindMany finds multiple documents in the specified collection
func (mc *MongoClient) FindMany(collectionName string, filter interface{}) (*mongo.Cursor, error) {
	collection := mc.getCollection(collectionName)
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %v", err)
	}
	return cursor, nil
}

// InsertMany inserts multiple documents into the specified collection
func (mc *MongoClient) UpdateOne(collectionName string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	collection := mc.getCollection(collectionName)
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %v", err)
	}
	return result, nil
}

// UpdateMany updates multiple documents in the specified collection
func (mc *MongoClient) UpdateMany(collectionName string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	collection := mc.getCollection(collectionName)
	result, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update documents: %v", err)
	}
	return result, nil
}

// DeleteOne deletes a single document from the specified collection
func (mc *MongoClient) DeleteOne(collectionName string, filter interface{}) (*mongo.DeleteResult, error) {
	collection := mc.getCollection(collectionName)
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete document: %v", err)
	}
	return result, nil
}

// DeleteMany deletes multiple documents from the specified collection
func (mc *MongoClient) DeleteMany(collectionName string, filter interface{}) (*mongo.DeleteResult, error) {
	collection := mc.getCollection(collectionName)
	result, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete documents: %v", err)
	}
	return result, nil
}
