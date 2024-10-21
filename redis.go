package godb

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var (
	redisClient *redis.Client
	redisOnce   sync.Once
)

// Initialize the Redis client with more options for timeout and retry.
func initialRedis() {
	redisHost := viper.GetString("REDIS_HOST")
	redisPassword := viper.GetString("REDIS_PASSWORD")
	redisIndexDB := viper.GetInt("REDIS_INDEX_DB")
	redisTimeout := viper.GetDuration("REDIS_TIMEOUT") // new configurable timeout
	if redisTimeout == 0 {
		redisTimeout = 10 * time.Second
	}

	if redisHost == "" {
		redisHost = "localhost:6379"
	}

	// Create the Redis client
	newClient := redis.NewClient(&redis.Options{
		Addr:         redisHost,
		Password:     redisPassword,
		DB:           redisIndexDB,
		DialTimeout:  redisTimeout,
		ReadTimeout:  redisTimeout,
		WriteTimeout: redisTimeout,
		PoolSize:     viper.GetInt("REDIS_POOL_SIZE"),   // Configurable pool size
		MaxRetries:   viper.GetInt("REDIS_MAX_RETRIES"), // Configurable retries
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()

	_, err := newClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	redisClient = newClient
	log.Println("Connected to Redis successfully")
}

// GetRedis returns the Redis client, initializing it if necessary.
func GetRedis() *redis.Client {
	redisOnce.Do(func() {
		initialRedis()
	})
	if redisClient == nil {
		log.Fatal("Redis client is nil: Redis uninitialized")
	}
	return redisClient
}

// ResetRedisCache flushes the entire Redis database.
func ResetRedisCache(client *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Set context timeout
	defer cancel()

	_, err := client.FlushAll(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to flush Redis cache: %v", err)
	}
	log.Println("Redis cache reset successfully")
	return nil
}

// CloseRedis cleans up the Redis connection when the app shuts down.
func CloseRedis() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}
