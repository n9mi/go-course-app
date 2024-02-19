package usecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/n9mi/go-course-app/internal/model"
	"github.com/n9mi/go-course-app/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type UserUseCase struct {
	ViperConfig    *viper.Viper
	DB             *gorm.DB
	Validate       *validator.Validate
	Log            *logrus.Logger
	UserRepository *repository.UserRepository
	RoleRepository *repository.RoleRepository
}

func NewUserUseCase(viperConfig *viper.Viper, db *gorm.DB, validate *validator.Validate, log *logrus.Logger,
	userRepository *repository.UserRepository, roleRepository *repository.RoleRepository) *UserUseCase {
	return &UserUseCase{
		ViperConfig:    viperConfig,
		DB:             db,
		Validate:       validate,
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
