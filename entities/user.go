package entities

import "github.com/keyloom/web-api/core"

type User struct {
	core.Entity `bson:",inline"`
	Username    string `bson:"username"`
	Email       string `bson:"email"`
	Password    string `bson:"password"`
}

var _ core.IEntity = (*User)(nil)

func (u *User) CreateNewEntity() *core.Entity {
	// Placeholder implementation
	return &core.Entity{}
}

func (u *User) LoadByID(id string) *core.Entity {
	// Placeholder implementation
	return &core.Entity{ID: id}
}

func (u *User) LoadByIDs(ids []string) []*core.Entity {
	// Placeholder implementation
	entities := make([]*core.Entity, len(ids))
	for i, id := range ids {
		entities[i] = &core.Entity{ID: id}
	}
	return entities
}

func (u *User) SaveOne(entity *core.Entity) error {
	// Placeholder implementation
	return nil
}

func (u *User) SaveMany(entities []*core.Entity) error {
	// Placeholder implementation
	return nil
}
