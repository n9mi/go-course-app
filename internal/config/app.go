package config

import (
	"github.com/casbin/casbin/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/n9mi/go-course-app/database/migration"
	"github.com/n9mi/go-course-app/database/seeder"
	"github.com/n9mi/go-course-app/internal/delivery/http/controller"
	"github.com/n9mi/go-course-app/internal/delivery/http/middleware"
	"github.com/n9mi/go-course-app/internal/delivery/http/route"
	"github.com/n9mi/go-course-app/internal/repository"
	"github.com/n9mi/go-course-app/internal/usecase"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ConfigBootstrap struct {
	ViperConfig *viper.Viper
	App         *fiber.App
	DB          *gorm.DB
	Validate    *validator.Validate
	RedisClient *redis.Client
	Log         *logrus.Logger
	Enforcer    *casbin.Enforcer
}

func Bootstrap(configBootstrap *ConfigBootstrap) {
	// Setup the repository
	repositorySetup := repository.Setup()

	// Setup usecase
	useCaseSetup := usecase.Setup(
		configBootstrap.ViperConfig,
		configBootstrap.DB,
		configBootstrap.Validate,
		configBootstrap.RedisClient,
		configBootstrap.Log,
		repositorySetup)

	// Setup controller
	controllerSetup := controller.Setup(useCaseSetup, configBootstrap.Log)

	// Setup middleware
	middlewareSetup := middleware.Setup(
		configBootstrap.ViperConfig,
		configBootstrap.Validate,
		configBootstrap.RedisClient,
		configBootstrap.Log,
		configBootstrap.Enforcer)

	// Setup Route
	routeConfig := route.RouteConfig{
		App:             configBootstrap.App,
		ControllerSetup: controllerSetup,
		MiddlewareSetup: middlewareSetup,
	}
	routeConfig.Setup()

	// Drop the database
	if err := migration.Drop(configBootstrap.DB); err != nil {
		configBootstrap.Log.Fatalf("Failed to drop the database: %v", err)
	}

	// Migrate the database
	if err := migration.Migrate(configBootstrap.DB); err != nil {
		configBootstrap.Log.Fatalf("Failed to migrate the database: %v", err)
	}

	// Seed the database
	if err := seeder.Seed(configBootstrap.DB, configBootstrap.RedisClient, repositorySetup); err != nil {
		configBootstrap.Log.Fatalf("Failed to seed the database: %v", err)
	}
}
