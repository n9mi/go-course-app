package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/n9mi/go-course-app/internal/model"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func GenerateAccessToken(log *logrus.Logger, viperConfig *viper.Viper, authData *model.UserAuthData) (string, int64, error) {
	key := viperConfig.GetString("ACCESS_TOKEN_KEY")
	expMinutes := viperConfig.GetInt("ACCESS_TOKEN_EXPIRE_MINUTES")

	return GenerateAuthToken(log, key, expMinutes, authData)
}

func GenerateRefreshToken(log *logrus.Logger, viperConfig *viper.Viper, authData *model.UserAuthData) (string, int64, error) {
	key := viperConfig.GetString("REFRESH_TOKEN_KEY")
	expMinutes := viperConfig.GetInt("REFRESH_TOKEN_EXPIRE_MINUTES")

	return GenerateAuthToken(log, key, expMinutes, authData)
}

func GenerateAuthToken(log *logrus.Logger, key string, expMinutes int, authData *model.UserAuthData) (string, int64, error) {
	timeDuration := time.Duration(expMinutes) * time.Minute
	expAtUnix := time.Now().Add(timeDuration).Unix()

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  expAtUnix,
		"data": *authData,
	})
	s, err := t.SignedString([]byte(key))
	if err != nil {
		log.Warnf("Failed to generate auth token : %+v", err)
		return "", 0, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return s, expAtUnix, nil
}

func ParseAccessToken(log *logrus.Logger, viperConfig *viper.Viper, token string) (*model.UserAuthData, error) {
	key := viperConfig.GetString("ACCESS_TOKEN_KEY")

	return ParseAuthToken(log, key, token)
}

func ParseRefreshToken(log *logrus.Logger, viperConfig *viper.Viper, token string) (*model.UserAuthData, error) {
	key := viperConfig.GetString("REFRESH_TOKEN_KEY")

	return ParseAuthToken(log, key, token)
}

func ParseAuthToken(log *logrus.Logger, key string, token string) (*model.UserAuthData, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		log.Warnf("Failed to parse auth token : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	// Parse auth data
	claims := t.Claims.(jwt.MapClaims)
	authDataMap, ok := claims["data"].(map[string]interface{})
	var authData model.UserAuthData

	authData.ID, ok = authDataMap["ID"].(string)
	if !ok {
		log.Warnf("Failed to parse auth data: Auth ID")
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	authData.Name, ok = authDataMap["Name"].(string)
	if !ok {
		log.Warnf("Failed to parse auth data: Auth Name")
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	authData.Email, ok = authDataMap["Email"].(string)
	if !ok {
		log.Warnf("Failed to parse auth data: Auth Email")
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	authRoles, ok := authDataMap["Roles"].([]interface{})
	if !ok {
		log.Warnf("Failed to parse auth data: Auth Roles")
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}
	authRolesData := make([]string, len(authRoles))
	for i := range authRoles {
		authRolesData[i], ok = authRoles[i].(string)
		if !ok {
			log.Warnf("Failed to parse auth data: Auth Roles %s", authRoles[i])
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
		}
	}
	authData.Roles = authRolesData

	return &authData, nil
}

func GetRefreshTokenRedisKey(userID string) string {
	return fmt.Sprintf("refreshKey:%s", userID)
}

func GetAccessTokenRedisKey(userID string) string {
	return fmt.Sprintf("accessKey:%s", userID)
}

func VerifyAccessToken(ctx context.Context, viperConfig *viper.Viper, redisClient *redis.Client,
	log *logrus.Logger, token string) (*model.UserAuthData, error) {
	// Parse access token
	tokenData, err := ParseAccessToken(log, viperConfig, token)
	if err != nil {
		return nil, err
	}

	// If token successfully parsed, then check get token from Redis
	tokenRedis, err := redisClient.Get(ctx, GetAccessTokenRedisKey(tokenData.ID)).Result()
	if err != nil {
		log.Warnf("Failed to get auth token from Redis : %v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	// Check if token that given from request is same with token that stored in Redis
	if tokenRedis != token {
		log.Warnf("Invalid auth token")
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	return tokenData, nil
}
