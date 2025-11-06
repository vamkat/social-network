// cmd/auth-service/main.go
package main

import (
	"log"

	"platform.zone01.gr/git/kvamvasa/social-network/pkg/config"
	"platform.zone01.gr/git/kvamvasa/social-network/pkg/db"
)

func main() {
	// 1️ Load environment config for Auth service
	cfg := config.Load(".env.auth")

	// 2️ Build DB configs
	pgCfg := db.PostgresConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPass,
		DBName:   cfg.DBName,
		SSLMode:  "disable",
	}

	rdCfg := db.RedisConfig{
		Addr: cfg.RedisUrl,
		DB:   0, // default Redis DB index
	}

	// 3️ Connect to Postgres
	pg, err := db.ConnectPostgres(pgCfg)
	if err != nil {
		log.Fatalf("❌ Postgres connection failed: %v", err)
	}
	defer pg.Close()

	// 4️ Run migrations for Auth service
	migrationPath := "pkg/db/migrations/auth"
	db.RunMigrations(cfg.DBUrl, migrationPath)

	// 5️ Connect to Redis
	rdb := db.ConnectRedis(rdCfg)
	defer rdb.Close()

	log.Println("🚀 Auth service database initialization completed successfully!")
}
