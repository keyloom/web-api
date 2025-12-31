package entities

import (
	"context"
	"time"

	"github.com/keyloom/web-api/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Migration struct {
	core.Entity `bson:",inline" json:",inline"`
	Changes     []string `bson:"changes" json:"changes"`
}

func (m *Migration) CollectionName() string {
	return "migrations"
}

func (m *Migration) CreateNew() *Migration {
	return &Migration{
		Entity: core.Entity{
			ID:        primitive.NilObjectID,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
		Changes: []string{},
	}
}

func (m *Migration) Save() error {
	client := core.NewMongoClient()
	if m.ID == primitive.NilObjectID {
		m.ID = primitive.NewObjectID()
		m.CreatedAt = time.Now().Unix()
		m.UpdatedAt = time.Now().Unix()
		_, err := client.InsertOne(m.CollectionName(), m)
		return err
	} else {
		m.UpdatedAt = time.Now().Unix()
		_, err := client.UpdateOne(m.CollectionName(), m.ID, m)
		return err
	}
}

func (m *Migration) GetLatest() (*Migration, error) {
	client := core.NewMongoClient()
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := client.FindMany(m.CollectionName(), bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}
	var migrations []Migration
	if err = cursor.All(context.TODO(), &migrations); err != nil {
		return nil, err
	}
	if len(migrations) == 0 {
		return nil, nil
	}
	return &migrations[0], nil
}
