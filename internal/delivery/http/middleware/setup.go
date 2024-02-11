package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type MiddlewareSetup struct {
	AuthMiddleware func(c *fiber.Ctx) error
}

func Setup(viperConfig *viper.Viper, validate *validator.Validate, redisClient *redis.Client,
	log *logrus.Logger) *MiddlewareSetup {
	return &MiddlewareSetup{
		AuthMiddleware: NewAuthMiddleware(viperConfig, redisClient, validate, log),
	}
}
