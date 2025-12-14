package core

import "go.mongodb.org/mongo-driver/bson/primitive"

type Entity struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt int64              `bson:"created_at"`
	UpdatedAt int64              `bson:"updated_at"`
}

type IEntity[T any] interface {
	CollectionName() string
	CreateNew() *T
	LoadByID(id string) *T
	LoadByIDs(ids []string) []*T
	Save() error
	Delete() error
}
