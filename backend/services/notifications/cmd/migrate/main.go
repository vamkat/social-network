package main

import (
	"log"
	"os"

	"social-network/shared/go/db"
)

func main() {
	log.Println("Running database migrations...")
	dbUrl := os.Getenv("DATABASE_URL")
	if err := db.RunMigrations(dbUrl, "./migrations"); err != nil {
		log.Fatal("Migration failed", err)
	}

	log.Println("Migrations completed successfully.")
	os.Exit(0)
}
