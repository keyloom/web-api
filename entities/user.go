package entities

import (
	"context"
	"time"

	"github.com/keyloom/web-api/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	core.Entity `json:",inline" bson:",inline"`
	Email       string `json:"email" bson:"email"`
	Password    string `json:"password" bson:"password,containsany=uppercase,containsany=lowercase,containsany=numeric,min=8"`
}

var _ core.IEntity[User] = (*User)(nil)

func (u *User) CollectionName() string {
	return "users"
}

func (u *User) CreateNew() *User {
	return &User{
		Entity: core.Entity{
			ID:        primitive.NilObjectID,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
	}
}

func (u *User) LoadByID(id string) *User {
	client := core.NewMongoClient()
	result := client.FindOne(u.CollectionName(), bson.M{"_id": id})
	if result.Err() != nil {
		return nil
	}
	var user User
	err := result.Decode(&user)
	if err != nil {
		return nil
	}
	return &user
}

// Loads multiple users by their IDs
func (u *User) LoadByIDs(ids []string) []*User {
	client := core.NewMongoClient()
	cursor, err := client.FindMany(u.CollectionName(), bson.M{
		"_id": bson.M{"$in": ids},
	})
	if err != nil {
		return nil
	}
	defer cursor.Close(context.TODO())

	var users []*User
	for cursor.Next(context.TODO()) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			continue
		}
		users = append(users, &user)
	}
	return users
}

// Saves the user to the database
func (u *User) Save() error {
	client := core.NewMongoClient()
	if u.ID != primitive.NilObjectID {
		_, err := client.UpdateOne(u.CollectionName(), bson.M{"_id": u.ID}, bson.M{"$set": u})
		return err
	} else {
		u.ID = primitive.NewObjectID()
		_, err := client.InsertOne(u.CollectionName(), u)
		return err
	}
}

// Deletes the user from the database
func (u *User) Delete() error {
	client := core.NewMongoClient()
	_, err := client.DeleteOne(u.CollectionName(), bson.M{"_id": u.ID})
	return err
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
