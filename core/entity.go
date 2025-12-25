package core

import "go.mongodb.org/mongo-driver/bson/primitive"

type Entity struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt int64              `bson:"created_at" json:"created_at"`
	UpdatedAt int64              `bson:"updated_at" json:"updated_at"`
}

type IEntity[T any] interface {
	CollectionName() string
	CreateNew() *T
	LoadByID(id string) *T
	LoadByIDs(ids []string) []*T
	LoadAll(top, page int) []*T
	Save() error
	Delete() error
}
