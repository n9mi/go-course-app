package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/n9mi/go-course-app/internal/entity"
	"github.com/n9mi/go-course-app/internal/helper"
	"github.com/n9mi/go-course-app/internal/model"
	"github.com/n9mi/go-course-app/internal/repository"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type AuthUseCase struct {
	ViperConfig    *viper.Viper
	DB             *gorm.DB
	Validate       *validator.Validate
	RedisClient    *redis.Client
	Log            *logrus.Logger
	UserRepository *repository.UserRepository
	RoleRepository *repository.RoleRepository
}

func NewAuthUseCase(viperConfig *viper.Viper, db *gorm.DB, validate *validator.Validate, redisClient *redis.Client,
	log *logrus.Logger, userRepository *repository.UserRepository, roleRepository *repository.RoleRepository) *AuthUseCase {
	return &AuthUseCase{
		ViperConfig:    viperConfig,
		DB:             db,
		Validate:       validate,
		RedisClient:    redisClient,
		Log:            log,
		UserRepository: userRepository,
		RoleRepository: roleRepository,
	}
}

func (u *AuthUseCase) Register(ctx context.Context, request *model.RegisterRequest) error {
	// Validate request
	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body : %+v", err)
		return err
	}

	tx := u.DB.WithContext(ctx).Begin()

	// Check if email are exists
	var userFound entity.User
	if u.UserRepository.FindByEmail(tx, &userFound, request.Email); len(userFound.ID) > 0 {
		u.Log.Warn("User already exists")
		return fiber.NewError(fiber.StatusConflict, "User already exists")
	}

	tx = u.DB.WithContext(ctx).Begin()

	// If email isn't exists yet, create user
	userPwd, err := helper.GeneratePassword(request.Password)
	if err != nil {
		u.Log.Warnf("Failed to generate password : %+v", err)
		return fiber.ErrInternalServerError
	}
	newUser := entity.User{
		ID:       "USR_" + helper.GenerateRandomString(16),
		Name:     request.Name,
		Email:    strings.ToLower(request.Email),
		Password: userPwd,
	}
	if err := u.UserRepository.Repository.Save(tx, &newUser); err != nil {
		if err := tx.Rollback().Error; err != nil {
			u.Log.Warnf("Failed to rollback users table : %+v", err)
			return fiber.ErrInternalServerError
		}
		u.Log.Warnf("Failed to create user : %+v", err)
		return fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed to commit users table : %+v", err)
		return fiber.ErrInternalServerError
	}

	// Attach "user" role to newly created user
	tx = u.DB.WithContext(ctx).Begin()
	var userRole entity.Role
	if err := u.RoleRepository.Repository.FindByID(tx, &userRole, "user"); err != nil {
		u.Log.Warnf("Failed to get 'user' role : %+v", err)
		return fiber.ErrInternalServerError
	}

	tx = u.DB.WithContext(ctx).Begin()
	if err := u.UserRepository.AssignRoles(tx, &newUser, []entity.Role{userRole}); err != nil {
		u.Log.Warnf("Failed to assign 'user' role to user with ID %s : %+v", newUser.ID, err)
		if err := tx.Rollback().Error; err != nil {
			u.Log.Warnf("Failed to rollback user_roles table : %+v", err)
			return fiber.ErrInternalServerError
		}
		return fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed to commit user_roles table : %+v", err)
		u.Log.Warnf("Failed to assign 'user' role to user with ID %s : %+v", newUser.ID, err)
		return fiber.ErrInternalServerError
	}

	return nil
}

func (u *AuthUseCase) Login(ctx context.Context, request *model.LoginRequest) (*model.TokenResponse, error) {
	// Validate request
	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body : %+v", err)
		return nil, err
	}

	tx := u.DB.WithContext(ctx).Begin()

	// Get user by email
	user := new(entity.User)
	if err := u.UserRepository.FindByEmail(tx, user, request.Email); err != nil || len(user.ID) == 0 {
		u.Log.Warnf("Failed to find user : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "User not found")
	}

	// If user is existed, check if password match
	if !helper.CompareHashAndPlainPassword(user.Password, request.Password) {
		u.Log.Warnf("Failed to authenticate the user : wrong password")
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Password doesn't match")
	}

	tx = u.DB.WithContext(ctx).Begin()

	userRole, err := u.UserRepository.GetRoles(tx, user)
	if err != nil {
		u.Log.Warnf("Failed to authenticate the user : role not found")
		return nil, fiber.NewError(fiber.StatusForbidden, "Forbidden")
	}
	var roles []string
	for _, r := range userRole {
		roles = append(roles, r.ID)
	}
	authData := model.UserAuthData{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Roles: roles,
	}

	// If authentication successful, generate new token
	response := new(model.TokenResponse)
	response.RefreshToken, response.RefreshExpAt, err = helper.GenerateRefreshToken(u.Log, u.ViperConfig, &authData)
	if err != nil {
		u.Log.Warnf("Failed to generate refresh token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	response.RefreshTokenName = u.ViperConfig.GetString("REFRESH_TOKEN_NAME")

	response.AccessToken, response.AccessExpAt, err = helper.GenerateAccessToken(u.Log, u.ViperConfig, &authData)
	if err != nil {
		u.Log.Warnf("Failed to generate access token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Store both on Redis
	refreshExpDur := response.RefreshExpAt - time.Now().Unix()
	if err := u.RedisClient.SetEx(ctx, helper.GetRefreshTokenRedisKey(authData.ID), response.RefreshToken,
		time.Duration(refreshExpDur)*time.Second).Err(); err != nil {
		u.Log.Warnf("Failed to save refresh token on Redis : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	accessExpDur := response.AccessExpAt - time.Now().Unix()
	if err := u.RedisClient.SetEx(ctx, helper.GetAccessTokenRedisKey(authData.ID), response.AccessToken,
		time.Duration(accessExpDur)*time.Second).Err(); err != nil {
		u.Log.Warnf("Failed to save access token on Redis : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return response, nil
}
