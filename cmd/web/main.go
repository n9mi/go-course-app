package main

import (
	"fmt"

	"github.com/n9mi/go-course-app/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	app := config.NewFiber(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	redisClient := config.NewRedisClient(viperConfig)
	validate := config.NewValidator(viperConfig)
	enforcer := config.NewEnforcer(log)

	configBootstrap := &config.ConfigBootstrap{
		ViperConfig: viperConfig,
		App:         app,
		DB:          db,
		Validate:    validate,
		RedisClient: redisClient,
		Log:         log,
		Enforcer:    enforcer,
	}
	config.Bootstrap(configBootstrap)

	webPort := viperConfig.GetInt("APP_PORT")
	if err := app.Listen(fmt.Sprintf(":%d", webPort)); err != nil {
		log.Fatalf("Failed to start the server : %+v", err)
	}
}
