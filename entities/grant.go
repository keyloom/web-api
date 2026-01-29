package entities

import (
	"context"
	"time"

	"github.com/keyloom/web-api/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Grant struct {
	core.Entity   `bson:",inline"`
	Scopes        []string           `bson:"scopes" json:"scopes"`
	UserID        primitive.ObjectID `bson:"userId" json:"-"`
	User          *User              `bson:"-" json:"user,omitempty"`
	ApplicationID primitive.ObjectID `bson:"applicationId" json:"-"`
	Application   *Application       `bson:"-" json:"application,omitempty"`
}

// CollectionName implements core.IEntity.
func (g *Grant) CollectionName() string {
	return "grants"
}

// CreateNew implements core.IEntity.
func (g *Grant) CreateNew() *Grant {
	return &Grant{
		Entity: core.Entity{
			ID:        primitive.NilObjectID,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
	}
}

// Delete implements core.IEntity.
func (g *Grant) Delete() error {
	client := core.NewMongoClient()
	_, err := client.DeleteOne(g.CollectionName(), bson.M{"_id": g.ID})
	return err
}

// LoadAll implements core.IEntity.
func (g *Grant) LoadAll(top int, page int) []*Grant {
	client := core.NewMongoClient()
	skip := (page - 1) * top
	findOptions := options.Find()
	findOptions.SetLimit(int64(top))
	findOptions.SetSkip(int64(skip))

	cursor, err := client.FindMany(g.CollectionName(), bson.D{}, findOptions)
	if err != nil {
		return nil
	}
	defer cursor.Close(context.TODO())

	var grants []*Grant
	for cursor.Next(context.TODO()) {
		var grant Grant
		err := cursor.Decode(&grant)
		if err != nil {
			continue
		}
		grants = append(grants, &grant)
	}

	return grants
}

// LoadByID implements core.IEntity.
func (g *Grant) LoadByID(id string) *Grant {
	client := core.NewMongoClient()
	filter := bson.M{"_id": id}
	result := client.FindOne(g.CollectionName(), filter)
	if result == nil {
		return nil
	}

	var grant Grant
	err := result.Decode(&grant)
	if err != nil {
		return nil
	}

	userEntity := &User{}
	user := userEntity.LoadByID(grant.UserID.Hex())
	grant.User = user

	return &grant
}

// LoadByIDs implements core.IEntity.
func (g *Grant) LoadByIDs(ids []string) []*Grant {
	client := core.NewMongoClient()
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := client.FindMany(g.CollectionName(), filter, nil)
	if err != nil {
		return nil
	}
	defer cursor.Close(context.TODO())

	var grants []*Grant
	for cursor.Next(context.TODO()) {
		var grant Grant
		err := cursor.Decode(&grant)
		if err != nil {
			continue
		}

		userEntity := &User{}
		user := userEntity.LoadByID(grant.UserID.Hex())
		grant.User = user

		grants = append(grants, &grant)
	}

	return grants
}

// Save implements core.IEntity.
func (g *Grant) Save() error {
	client := core.NewMongoClient()
	if g.ID != primitive.NilObjectID {
		g.UpdatedAt = time.Now().Unix()
		_, err := client.UpdateOne(g.CollectionName(), bson.M{"_id": g.ID}, g)
		return err
	} else {
		g.ID = primitive.NewObjectID()
		_, err := client.InsertOne(g.CollectionName(), g)
		return err
	}
}

func (g *Grant) CreateDefaultGrant(migration *Migration) error {
	adminUserConfig, err := (&core.EnvManager{}).GetAdminUserConfig()
	if err != nil {
		return err
	}
	userEntity := &User{}
	adminUser := userEntity.LoadByEmail(adminUserConfig.Email)
	if adminUser == nil {
		return nil
	}

	applicationEntity := &Application{}
	defaultApp := applicationEntity.LoadByName("keyloom-frontend")
	if defaultApp == nil {
		return nil
	}

	// Check if default grant exists
	client := core.NewMongoClient()
	result := client.FindOne(g.CollectionName(), bson.M{
		"userId":        adminUser.ID.Hex(),
		"applicationId": defaultApp.ID.Hex(),
	})
	if result.Err() == nil {
		// Grant already exists
		return nil
	}

	// Create default grant
	defaultGrant := &Grant{
		Scopes:        []string{"read", "write", "delete"},
		UserID:        adminUser.ID,
		ApplicationID: defaultApp.ID,
	}
	err = defaultGrant.Save()
	if err != nil {
		return err
	}
	// Update migration record
	migration.Changes = append(migration.Changes, core.MigrationChangeCreateDefaultGrant)
	return nil
}

var _ core.IEntity[Grant] = (*Grant)(nil)
