package usecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/n9mi/go-course-app/internal/model"
	"github.com/n9mi/go-course-app/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CategoryUseCase struct {
	DB                 *gorm.DB
	Validate           *validator.Validate
	Log                *logrus.Logger
	CategoryRepository *repository.CategoryRepository
}

func NewCategoryUseCase(db *gorm.DB, validate *validator.Validate, log *logrus.Logger,
	categoryRepository *repository.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{
		DB:                 db,
		Validate:           validate,
		Log:                log,
		CategoryRepository: categoryRepository,
	}
}

func (u *CategoryUseCase) List(ctx context.Context, request *model.CategoryListRequest) ([]model.CategoryResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()

	categories, err := u.CategoryRepository.List(tx, request)
	if err != nil {
		u.Log.Warnf("Failed to get categories : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return categories, nil
}
