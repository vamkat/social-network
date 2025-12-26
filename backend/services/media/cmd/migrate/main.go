package main

import (
	"log"
	"os"

	"social-network/shared/go/db"
)

func main() {
	log.Println("Running database migrations...")

	// if err := db.RunMigrations(os.Getenv("DATABASE_URL"), "./migrations"); err != nil {
	if err := db.RunMigrations("postgres://postgres:secret@localhost:5437/social_media?sslmode=disable", "services/media/internal/db/migrations"); err != nil {
		log.Fatal("migration failed", err)
	}

	log.Println("Migrations completed successfully.")
	os.Exit(0)
}
