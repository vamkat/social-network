package redis_connector

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrNoConnection = errors.New("can't find redis")
	ErrFailedTest   = errors.New("failed special connection test")
)

type redisClient struct {
	client *redis.Client
}

func NewRedisClient(addr string, password string, db int) redisClient {
	redisClient := redisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}

	return redisClient
}

func (c *redisClient) SetStr(ctx context.Context, key string, value string, exp time.Duration) error {
	err := c.client.Set(ctx, key, value, exp).Err()
	return err
}

func (c *redisClient) GetStr(ctx context.Context, key string) (any, error) {
	value, err := c.client.Get(ctx, key).Result()
	return value, err
}

// SetObj marshals a Go value to JSON and stores it under `key` with expiration `exp`.
func (c *redisClient) SetObj(ctx context.Context, key string, value any, exp time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, b, exp).Err()
}

// GetObj retrieves the JSON stored at `key` and unmarshals it into `dest`.
// `dest` must be a pointer to the value to populate.
func (c *redisClient) GetObj(ctx context.Context, key string, dest any) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (c *redisClient) Del(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	return err
}

func (c *redisClient) TestRedisConnection() error {
	ctx := context.Background()
	ping, err := c.client.Ping(ctx).Result()
	if err != nil || ping != "PONG" {
		return errors.Join(ErrNoConnection, err)
	}

	err = c.SetStr(ctx, "test_key123", "value", time.Second)
	if err != nil {
		return ErrFailedTest
	}

	val, err := c.GetStr(ctx, "test_key123")
	valStr, ok := val.(string)
	if err != nil || !ok || valStr != "value" {
		return ErrFailedTest
	}

	err = c.Del(ctx, "test_key123")
	if err != nil {
		return ErrFailedTest
	}

	_, err = c.GetStr(ctx, "test_key123")
	if err != redis.Nil {
		return ErrFailedTest
	}
	return nil
}
