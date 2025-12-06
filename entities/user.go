package entities

import (
	"github.com/keyloom/web-api/core"
)

type User struct {
	core.Entity `bson:",inline"`
	Email       string `bson:"email,email"`
	Password    string `bson:"password,containsany=uppercase,containsany=lowercase,containsany=numeric,min=8"`
}

var _ core.IEntity = (*User)(nil)

func (u *User) CollectionName() string {
	return "users"
}

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

// Hashes and sets the user's password
func (u *User) SetPassword(password string) error {
	// Hash the password before storing it
	hasher := core.Hasher{}
	hashedPassword, err := hasher.Hash(password)
	if err != nil {
		return err
	}

	// Set the hashed password
	u.Password = hashedPassword
	return nil
}

// Compares the given password with the stored hashed password
func (u *User) CheckPassword(password string) bool {
	hasher := core.Hasher{}
	return hasher.Compare(u.Password, password)
}

// Sets the user's email
func (u *User) SetEmail(email string) {
	u.Email = email
}
