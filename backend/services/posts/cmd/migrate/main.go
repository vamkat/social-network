package main

import (
	"log"
	"os"
	"time"

	"social-network/shared/go/db"
)

func main() {
	log.Println("Running database migrations...")
	dbUrl := os.Getenv("DATABASE_URL")
	for range 10 {
		if err := db.RunMigrations(dbUrl, "./migrations"); err != nil {
			log.Println("Migration failed, retrying in 2s:", err)
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	log.Println("Migrations completed successfully.")
	os.Exit(0)
}
