package main

import (
	"fmt"
	"log"

	"product-api-assessment/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrate(cfg *config.Config, direction string) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name,
	)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	switch direction {
	case "up":
		migrateUp(m)
	case "down":
		migrateDown(m)
	default:
		log.Fatalf("Unknown migrate direction: %s (use 'up' or 'down')", direction)
	}
}

func migrateUp(m *migrate.Migrate) {
	for {
		version, dirty, _ := m.Version()
		if err := m.Steps(1); err != nil {
			if err == migrate.ErrNoChange || err.Error() == "file does not exist" {
				break
			}
			log.Fatalf("Migration up failed: %v", err)
		}
		newVersion, _, _ := m.Version()
		if dirty {
			log.Printf("Warning: version %d was dirty", version)
		}
		log.Printf("Migrated up: version %d → %d", version, newVersion)
	}
	log.Println("Migration up completed successfully")
}

func migrateDown(m *migrate.Migrate) {
	for {
		version, _, err := m.Version()
		if err != nil {
			break
		}
		log.Printf("Migrating down: version %d", version)
		if err := m.Steps(-1); err != nil {
			if err == migrate.ErrNoChange {
				break
			}
			log.Fatalf("Migration down failed: %v", err)
		}
	}
	log.Println("Migration down completed successfully")
}
