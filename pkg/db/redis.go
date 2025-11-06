// pkg/db/redis.go
package db

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// ConnectRedis creates a new Redis client using config struct
func ConnectRedis(cfg RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		DB:           cfg.DB,
		Password:     cfg.Password,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis (%s): %v", cfg.Addr, err)
	}

	log.Printf("✅ Connected to Redis at %s (DB %d)", cfg.Addr, cfg.DB)
	return rdb
}
