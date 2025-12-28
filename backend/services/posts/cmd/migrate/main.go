package main

import (
	"log"
	"os"

	"social-network/shared/go/db"
	tele "social-network/shared/go/telemetry"
)

func main() {
	log.Println("Running database migrations...")
	dbUrl := os.Getenv("DATABASE_URL")
	if err := db.RunMigrations(dbUrl, os.Getenv("MIGRATE_PATH")); err != nil {
		tele.Fatalf("Migration failed: %v", err)
	}

	tele.Info(nil, "Migrations completed successfully.")
	os.Exit(0)
}
