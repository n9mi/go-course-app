package usecase

import (
	"context"

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

type UserUseCase struct {
	ViperConfig    *viper.Viper
	DB             *gorm.DB
	Validate       *validator.Validate
	RedisClient    *redis.Client
	Log            *logrus.Logger
	UserRepository *repository.UserRepository
	RoleRepository *repository.RoleRepository
}

func NewUserUseCase(viperConfig *viper.Viper, db *gorm.DB, validate *validator.Validate, redisClient *redis.Client,
	log *logrus.Logger, userRepository *repository.UserRepository, roleRepository *repository.RoleRepository) *UserUseCase {
	return &UserUseCase{
		ViperConfig:    viperConfig,
		DB:             db,
		Validate:       validate,
		RedisClient:    redisClient,
		Log:            log,
		UserRepository: userRepository,
		RoleRepository: roleRepository,
	}
}

func (u *UserUseCase) List(ctx context.Context, request *model.UserListRequest) ([]model.UserListResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()

	userEntities, err := u.UserRepository.List(tx, request)
	if err != nil {
		u.Log.Warnf("Failed to get users : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	users := make([]model.UserListResponse, len(userEntities))
	for i, uE := range userEntities {
		users[i] = model.UserListResponse{
			ID:       uE.ID,
			Name:     uE.Name,
			Email:    uE.Email,
			JoinedAt: uE.CreatedAt,
		}
		for _, uR := range userEntities[i].Roles {
			users[i].Roles = append(users[i].Roles, model.RoleListResponse{
				Name:    uR.DisplayName,
				AddedAt: uR.CreatedAt,
			})
		}
	}

	return users, nil
}

func (u *UserUseCase) Delete(ctx context.Context, request *model.UserDeleteRequest) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Failed to validate user : %+v", err)
		return err
	}

	user := new(entity.User)
	if err := u.UserRepository.Repository.FindByID(tx, user, request.ID); err != nil {
		u.Log.Warnf("User with ID %s is not found : %+v", request.ID, err)
		return fiber.ErrNotFound
	}

	if err := u.RedisClient.Del(ctx, helper.GetAccessTokenRedisKey(user.ID),
		helper.GetRefreshTokenRedisKey(user.ID)).Err(); err != nil {
		u.Log.Warnf("Failed to auth token from user with ID %s : %+v", request.ID, err)
		return fiber.ErrInternalServerError
	}

	if err := u.UserRepository.Delete(tx, user); err != nil {
		u.Log.Warnf("Failed to delete user with ID %s : %+v", request.ID, err)
		return fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed to commit into 'users' table : %+v", err)
		return fiber.ErrInternalServerError
	}

	return nil
}
