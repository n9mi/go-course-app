package helper

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/n9mi/go-course-app/internal/model"
	"github.com/spf13/viper"
)

func GenerateAccessToken(viperConfig *viper.Viper, authData *model.UserAuthData) (string, int64, error) {
	key := viperConfig.GetString("ACCESS_TOKEN_KEY")
	expMinutes := viperConfig.GetInt("ACCESS_TOKEN_EXPIRE_MINUTES")

	return GenerateAuthToken(key, expMinutes, authData)
}

func GenerateRefreshToken(viperConfig *viper.Viper, authData *model.UserAuthData) (string, int64, error) {
	key := viperConfig.GetString("REFRESH_TOKEN_KEY")
	expMinutes := viperConfig.GetInt("REFRESH_TOKEN_EXPIRE_MINUTES")

	return GenerateAuthToken(key, expMinutes, authData)
}

func GenerateAuthToken(key string, expMinutes int, authData *model.UserAuthData) (string, int64, error) {
	timeDuration := time.Duration(expMinutes) * time.Minute
	expAtUnix := time.Now().Add(timeDuration).Unix()

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  expAtUnix,
		"data": *authData,
	})
	s, err := t.SignedString([]byte(key))
	if err != nil {
		return "", 0, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return s, expAtUnix, nil
}

func ParseAccessToken(viperConfig *viper.Viper, token string) (*model.UserAuthData, error) {
	key := viperConfig.GetString("ACCESS_TOKEN_KEY")
	expMinutes := viperConfig.GetInt("ACCESS_TOKEN_EXPIRE_MINUTES")

	return ParseAuthToken(key, expMinutes, token)
}

func ParseRefreshToken(viperConfig *viper.Viper, token string) (*model.UserAuthData, error) {
	key := viperConfig.GetString("REFRESH_TOKEN_KEY")
	expMinutes := viperConfig.GetInt("REFRESH_TOKEN_EXPIRE_MINUTES")

	return ParseAuthToken(key, expMinutes, token)
}

func ParseAuthToken(key string, expMinutes int, token string) (*model.UserAuthData, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	claims := t.Claims.(jwt.MapClaims)
	timeExp := int64(claims["exp"].(float64))

	// Check if token is expired
	if time.Now().Unix() > timeExp {
		return nil, fiber.NewError(http.StatusBadRequest, "expired token")
	}

	// Parse auth data
	data, ok := claims["data"].(model.UserAuthData)
	if !ok {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid token")
	}

	return &data, nil
}

func GetRefreshRedisKey(userID string) string {
	return fmt.Sprintf("refreshKey:%s", userID)
}

func GetAccessRedisKey(userID string) string {
	return fmt.Sprintf("accessKey:%s", userID)
}
