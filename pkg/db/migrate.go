package db

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbURL, migrationPath string) {
	m, err := migrate.New("file://"+migrationPath, dbURL)
	if err != nil {
		log.Fatalf("❌ Migration init failed: %v", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("⚠️ No new migrations to apply")
		} else {
			log.Fatalf("❌ Migration failed: %v", err)
		}
	} else {
		log.Println("✅ Migrations applied successfully")
	}
}
