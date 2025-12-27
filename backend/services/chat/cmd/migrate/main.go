package main

import (
	"log"
	"os"

	"social-network/shared/go/db"
)

func main() {
	log.Println("Running database migrations...")

	if err := db.RunMigrations(os.Getenv("DATABASE_URL"), os.Getenv("MIGRATE_PATH")); err != nil {
		log.Fatal("Migration failed", err)
	}

	log.Println("Migrations completed successfully.")
	os.Exit(0)
}
