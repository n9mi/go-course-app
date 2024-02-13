package usecase

import (
	"github.com/go-playground/validator/v10"
	"github.com/n9mi/go-course-app/internal/repository"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type UseCaseSetup struct {
	AuthUseCase     *AuthUseCase
	CategoryUseCase *CategoryUseCase
}

func Setup(viperConfig *viper.Viper, db *gorm.DB, validate *validator.Validate, redisClient *redis.Client,
	log *logrus.Logger, repositorySetup *repository.RepositorySetup) *UseCaseSetup {

	return &UseCaseSetup{
		AuthUseCase: NewAuthUseCase(viperConfig, db, validate, redisClient, log,
			repositorySetup.UserRepository, repositorySetup.RoleRepository),
		CategoryUseCase: NewCategoryUseCase(db, validate, log, repositorySetup.CategoryRepository),
	}
}
