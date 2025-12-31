package entities

import (
	"context"
	"time"

	"github.com/keyloom/web-api/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ClientSecret struct {
	Name      string `bson:"name" json:"name"`
	Value     string `bson:"value" json:"value"`
	ExpireAt  int64  `bson:"expire_at" json:"expire_at"`
	CreatedAt int64  `bson:"created_at" json:"created_at"`
}

type Application struct {
	core.Entity   `bson:",inline" json:",inline"`
	Name          string         `bson:"name" json:"name"`
	Description   string         `bson:"description" json:"description"`
	ClientID      string         `bson:"client_id" json:"client_id"`
	ClientSecrets []ClientSecret `bson:"client_secret" json:"client_secret"`
}

var _ core.IEntity[Application] = (*Application)(nil)

func (a *Application) CollectionName() string {
	return "applications"
}

func (a *Application) CreateNew() *Application {
	return &Application{
		Entity: core.Entity{
			ID:        primitive.NilObjectID,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
		ClientSecrets: []ClientSecret{},
	}
}

func (a *Application) LoadAll(top, page int) []*Application {
	client := core.NewMongoClient()
	skip := (page - 1) * top
	findOptions := options.Find()
	findOptions.SetLimit(int64(top))
	findOptions.SetSkip(int64(skip))

	cursor, err := client.FindMany(a.CollectionName(), bson.D{}, findOptions)
	if err != nil {
		return nil
	}
	defer cursor.Close(context.TODO())

	var applications []*Application
	for cursor.Next(context.TODO()) {
		var app Application
		err := cursor.Decode(&app)
		if err != nil {
			continue
		}
		applications = append(applications, &app)
	}
	return applications
}

func (a *Application) LoadByID(id string) *Application {
	client := core.NewMongoClient()
	result := client.FindOne(a.CollectionName(), map[string]interface{}{"_id": id})
	if result.Err() != nil {
		return nil
	}
	var application Application
	err := result.Decode(&application)
	if err != nil {
		return nil
	}
	return &application
}

// Loads multiple applications by their IDs
func (a *Application) LoadByIDs(ids []string) []*Application {
	client := core.NewMongoClient()
	cursor, err := client.FindMany(a.CollectionName(), bson.M{
		"_id": bson.M{"$in": ids},
	})
	if err != nil {
		return nil
	}
	defer cursor.Close(context.TODO())

	var applications []*Application
	for cursor.Next(context.TODO()) {
		var app Application
		err := cursor.Decode(&app)
		if err != nil {
			continue
		}
		applications = append(applications, &app)
	}
	return applications
}

func (a *Application) Save() error {
	client := core.NewMongoClient()
	if a.ID != primitive.NilObjectID {
		a.UpdatedAt = time.Now().Unix()
		_, err := client.UpdateOne(a.CollectionName(), bson.M{"_id": a.ID}, bson.M{"$set": a})
		return err
	} else {
		a.ID = primitive.NewObjectID()
		a.CreatedAt = time.Now().Unix()
		a.UpdatedAt = time.Now().Unix()
		_, err := client.InsertOne(a.CollectionName(), a)
		return err
	}
}

func (a *Application) Delete() error {
	client := core.NewMongoClient()
	_, err := client.DeleteOne(a.CollectionName(), bson.M{"_id": a.ID})
	return err
}
