package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis() error {
	// Specify the Redis address directly
	addr := "red-cqvpisjv2p9s739fmp50:6379" // Replace with your Redis address from Render

	Rdb = redis.NewClient(&redis.Options{
		Addr: addr, // Use the address directly
		DB:   0,    // default DB
	})

	// Test the connection to Redis
	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		return fmt.Errorf("could not connect to Redis: %v", err)
	}

	fmt.Println("Successfully connected to Redis")
	return nil
}
