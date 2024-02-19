package middleware

import (
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/n9mi/go-course-app/internal/helper"
	"github.com/n9mi/go-course-app/internal/model"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewAuthMiddleware(viperConfig *viper.Viper, redisClient *redis.Client, validate *validator.Validate,
	log *logrus.Logger, enforcer *casbin.Enforcer) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization", "")

		// Get token from Bearer
		if !strings.Contains(header, "Bearer ") {
			log.Warnf("Authorization token hasn't provided")
			return fiber.ErrUnauthorized
		}
		accessToken := strings.Replace(header, "Bearer ", "", -1)

		// Check if token is provided in header
		request := model.AuthToken{Token: accessToken}
		if err := validate.Struct(request); err != nil {
			log.Warnf("Authorization token hasn't provided")
			return fiber.ErrUnauthorized
		}

		// Check if token is valid and get the auth data from the provided token
		userAuthData, err := helper.VerifyAccessToken(c.UserContext(), viperConfig, redisClient, log, accessToken)
		if err != nil {
			return err
		}

		// Apply casbin enforcer, loop through roles
		roles := userAuthData.Roles
		isRoleExist := false
		for _, casbinSubject := range roles {
			enforce, err := enforcer.Enforce(casbinSubject, c.Path(), c.Method())
			if err != nil {
				log.Warnf("Failed to enforce path")
				return fiber.ErrInternalServerError
			}
			if enforce {
				isRoleExist = true
			}
		}
		if isRoleExist {
			// Store auth data in Locals
			c.Locals("AuthData", *userAuthData)

			return c.Next()
		}

		log.Warnf("Forbidden error, user role doesn't match")
		return fiber.ErrForbidden
	}
}
