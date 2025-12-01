package main

import (
	"log"
	"os"
	"time"

	"social-network/shared/go/db"
)

func main() {
	log.Println("Running database migrations...")
	for range 10 {
		if err := db.RunMigrations(os.Getenv("DATABASE_URL"), "./migrations"); err != nil {
			log.Println("Migration failed, retrying in 2s:", err)
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	log.Println("Migrations completed successfully.")
	os.Exit(0)
}
