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
	core.Entity       `bson:",inline" json:",inline"`
	Name              string               `bson:"name" json:"name"`
	Description       string               `bson:"description" json:"description"`
	ClientID          string               `bson:"client_id" json:"client_id"`
	ClientSecrets     []ClientSecret       `bson:"client_secret" json:"client_secret"`
	RedirectURIs      []string             `bson:"redirect_uris" json:"redirect_uris"`
	Scopes            []string             `bson:"scopes" json:"scopes"`
	ResourceServerIDs []primitive.ObjectID `bson:"resource_server_ids" json:"-"`
	ResourceServers   []*ResourceServer    `bson:"-" json:"resource_servers,omitempty"`
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
		ClientSecrets:     []ClientSecret{},
		RedirectURIs:      []string{},
		Scopes:            []string{},
		ResourceServerIDs: []primitive.ObjectID{},
		ResourceServers:   []*ResourceServer{},
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
		resourceServerEntity := &ResourceServer{}
		stringIds := make([]string, len(app.ResourceServerIDs))
		for i, id := range app.ResourceServerIDs {
			stringIds[i] = id.Hex()
		}
		resourceServers := resourceServerEntity.LoadByIDs(stringIds)
		app.ResourceServers = resourceServers
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
	resourceServerEntity := &ResourceServer{}
	stringIds := make([]string, len(application.ResourceServerIDs))
	for i, id := range application.ResourceServerIDs {
		stringIds[i] = id.Hex()
	}
	resourceServers := resourceServerEntity.LoadByIDs(stringIds)
	application.ResourceServers = resourceServers
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
		resourceServerEntity := &ResourceServer{}
		stringIds := make([]string, len(app.ResourceServerIDs))
		for i, id := range app.ResourceServerIDs {
			stringIds[i] = id.Hex()
		}
		resourceServers := resourceServerEntity.LoadByIDs(stringIds)
		app.ResourceServers = resourceServers
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

func (a *Application) LoadByName(name string) *Application {
	client := core.NewMongoClient()
	result := client.FindOne(a.CollectionName(), bson.M{"name": name})
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

func (a *Application) CreateDefaultApplication(migration *Migration) error {
	client := core.NewMongoClient()
	// Check if default application exists
	result := client.FindOne(a.CollectionName(), bson.M{"name": "keyloom-frontend"})
	if result.Err() == nil {
		// Default application already exists
		return nil
	}
	// Create default application
	defaultApp := a.CreateNew()
	defaultApp.Name = "keyloom-frontend"
	defaultApp.Description = "Default Keyloom Frontend Application"
	defaultApp.ClientID = "keyloom-frontend-client-id"
	defaultApp.RedirectURIs = []string{
		"http://localhost:3000/callback",
		"http://localhost:3000/redirect",
	}
	defaultApp.Scopes = []string{
		"keyloom:view:resource-servers",
		"keyloom:manage:resource-servers",
		"keyloom:view:applications",
		"keyloom:manage:applications",
		"keyloom:view:users",
		"keyloom:manage:users",
		"keyloom:view:grants",
		"keyloom:manage:grants",
	}

	err := defaultApp.Save()
	if err != nil {
		return err
	}

	// Record migration change
	migration.Changes = append(migration.Changes, core.MigrationChangeCreateDefaultApplication)
	return nil
}
