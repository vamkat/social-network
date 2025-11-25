package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Run migrations from a given path
func RunMigrations(dbUrl string, migrationsPath string) error {
	sqlDB, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return fmt.Errorf("failed to open DB for migrations: %w", err)
	}
	defer sqlDB.Close()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres", driver,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("✅ Migrations applied successfully")
	return nil
}

// func RunMigrations(cfg Config, migrationsPath string) error {
// 	sqlDB, err := sql.Open("postgres", cfg.ConnString())
// 	if err != nil {
// 		return fmt.Errorf("failed to open DB for migrations: %w", err)
// 	}
// 	defer sqlDB.Close()

// 	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
// 	if err != nil {
// 		return fmt.Errorf("failed to create migrate driver: %w", err)
// 	}

// 	m, err := migrate.NewWithDatabaseInstance(
// 		fmt.Sprintf("file://%s", migrationsPath),
// 		"postgres", driver,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("failed to initialize migrate: %w", err)
// 	}

// 	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
// 		return fmt.Errorf("migration failed: %w", err)
// 	}

// 	log.Println("✅ Migrations applied successfully")
// 	return nil
// }
