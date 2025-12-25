package entities

import (
	"context"
	"time"

	"github.com/keyloom/web-api/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Audience struct {
	core.Entity `bson:",inline" json:",inline"`
	DisplayName string `bson:"display_name" json:"display_name"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
}

var _ core.IEntity[Audience] = (*Audience)(nil)

func (a *Audience) CollectionName() string {
	return "audiences"
}

func (a *Audience) CreateNew() *Audience {
	return &Audience{
		Entity: core.Entity{
			ID:        primitive.NilObjectID,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
	}
}

func (a *Audience) LoadAll(top, page int) []*Audience {
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

	var audiences []*Audience
	for cursor.Next(context.TODO()) {
		var audience Audience
		err := cursor.Decode(&audience)
		if err != nil {
			continue
		}
		audiences = append(audiences, &audience)
	}
	return audiences
}

func (a *Audience) LoadByID(id string) *Audience {
	client := core.NewMongoClient()
	result := client.FindOne(a.CollectionName(), map[string]interface{}{"_id": id})
	if result.Err() != nil {
		return nil
	}
	var audience Audience
	err := result.Decode(&audience)
	if err != nil {
		return nil
	}
	return &audience
}

// Loads multiple audiences by their IDs
func (a *Audience) LoadByIDs(ids []string) []*Audience {
	client := core.NewMongoClient()
	cursor, err := client.FindMany(a.CollectionName(), map[string]interface{}{
		"_id": map[string]interface{}{"$in": ids},
	})
	if err != nil {
		return nil
	}
	defer cursor.Close(context.TODO())

	var audiences []*Audience
	for cursor.Next(context.TODO()) {
		var audience Audience
		err := cursor.Decode(&audience)
		if err != nil {
			continue
		}
		audiences = append(audiences, &audience)
	}
	return audiences
}

func (a *Audience) Save() error {
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

func (a *Audience) Delete() error {
	client := core.NewMongoClient()
	_, err := client.DeleteOne(a.CollectionName(), bson.M{"_id": a.ID})
	return err
}
