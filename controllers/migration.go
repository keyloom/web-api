package controllers

import (
	"fmt"
	"slices"

	"github.com/keyloom/web-api/core"
	"github.com/keyloom/web-api/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MigrationController struct{}

func (mc *MigrationController) RunMigrations() {
	fmt.Println("")
	fmt.Println("[MIGRATIONS] Starting migrations...")
	// Create migration Object
	migration := &entities.Migration{}
	latestMigration, err := migration.GetLatest()
	if err != nil {
		return
	}
	if latestMigration == nil {
		latestMigration = migration.CreateNew()
	}

	fmt.Println("")
	if latestMigration.ID == primitive.NilObjectID {
		fmt.Println("[MIGRATIONS] No migrations found.")
		fmt.Println("[MIGRATIONS] New migration record created.")
	} else {
		fmt.Printf("[MIGRATIONS] Latest migration found: %s\n", latestMigration.ID.Hex())
	}

	// Create default admin user if not exists
	if !slices.Contains(latestMigration.Changes, core.MigrationChangeCreateDefaultAdminUser) {
		fmt.Println("[MIGRATIONS] No default admin user found. Creating one...")
		userEntity := &entities.User{}
		err := userEntity.CreateDefaultAdminUser(latestMigration)
		if err != nil {
			return
		}
		fmt.Println("[MIGRATIONS] Default admin user created.")
	}

	// Create default resource server if not exists
	if !slices.Contains(latestMigration.Changes, core.MigrationChangeCreateDefaultResourceServer) {
		fmt.Println("[MIGRATIONS] No default resource server found. Creating one...")
		resourceServerEntity := &entities.ResourceServer{}
		err := resourceServerEntity.CreateDefaultResourceServer(latestMigration)
		if err != nil {
			return
		}
		fmt.Println("[MIGRATIONS] Default resource server created.")
	}

	// Create default application if not exists
	if !slices.Contains(latestMigration.Changes, core.MigrationChangeCreateDefaultApplication) {
		fmt.Println("[MIGRATIONS] No default application found. Creating one...")
		applicationEntity := &entities.Application{}
		err := applicationEntity.CreateDefaultApplication(latestMigration)
		if err != nil {
			return
		}
		fmt.Println("[MIGRATIONS] Default application created.")
	}

	// Create default grant if not exists
	if !slices.Contains(latestMigration.Changes, core.MigrationChangeCreateDefaultGrant) {
		fmt.Println("[MIGRATIONS] No default grant found. Creating one...")
		grantEntity := &entities.Grant{}
		err := grantEntity.CreateDefaultGrant(latestMigration)
		if err != nil {
			return
		}
		fmt.Println("[MIGRATIONS] Default grant created.")
	}

	fmt.Println("")
	fmt.Println("[MIGRATIONS] Migrations completed.")
	// Save latest migration
	latestMigration.Save()
	fmt.Println("")
	fmt.Printf("[MIGRATIONS] Latest migration ID: %s\n", latestMigration.ID.Hex())
	fmt.Println("")
}
