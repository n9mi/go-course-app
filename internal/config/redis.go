package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// NewRedisClient returns redis client instance with configuration
// that provided in config.env
func NewRedisClient(config *viper.Viper) *redis.Client {
	// Getting config from viper
	address := config.GetString("REDIS_ADDRESS")
	port := config.GetString("REDIS_PORT")
	db := config.GetInt("REDIS_DB")
	password := config.GetString("REDIS_PASSWORD")

	// Creating new client
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", address, port),
		DB:       db,
		Password: password,
	})

	// Delete all redis value when app restarting
	client.FlushDB(context.Background())

	return client
}
