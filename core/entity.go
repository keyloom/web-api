package core

type Entity struct {
	ID        string `bson:"_id"`
	CreatedAt int64  `bson:"created_at"`
	UpdatedAt int64  `bson:"updated_at"`
}

type IEntity interface {
	CollectionName() string
	CreateNewEntity() *Entity
	LoadByID(id string) *Entity
	LoadByIDs(ids []string) []*Entity
	SaveOne(entity *Entity) error
	SaveMany(entities []*Entity) error
}
