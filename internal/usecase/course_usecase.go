package usecase

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/n9mi/go-course-app/internal/model"
	"github.com/n9mi/go-course-app/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CourseUseCase struct {
	DB                     *gorm.DB
	Validate               *validator.Validate
	Log                    *logrus.Logger
	Cloudinary             *cloudinary.Cloudinary
	CourseRepository       *repository.CourseRepository
	CourseMemberRepository *repository.CourseMemberRepository
}

func NewCourseUseCase(db *gorm.DB, validate *validator.Validate, log *logrus.Logger, cld *cloudinary.Cloudinary,
	courseRepository *repository.CourseRepository, courseMemberRepository *repository.CourseMemberRepository) *CourseUseCase {
	return &CourseUseCase{
		DB:                     db,
		Validate:               validate,
		Log:                    log,
		Cloudinary:             cld,
		CourseRepository:       courseRepository,
		CourseMemberRepository: courseMemberRepository,
	}
}

func (u *CourseUseCase) List(ctx context.Context, request *model.CourseListRequest) ([]model.CourseListResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Failed to validate category : %+v", err)
		return nil, err
	}

	response, err := u.CourseRepository.List(tx, request)
	if err != nil {
		u.Log.Warnf("Failed to get courses : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return response, nil
}

func (u *CourseUseCase) FindByID(ctx context.Context, request *model.CourseGetRequest) (*model.CourseResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()

	course := new(model.CourseResponse)
	if err := u.CourseRepository.FindByIDAndUserID(tx, course, request.ID, request.UserID); err != nil {
		u.Log.Warnf("Failed to find course with ID %s : %+v", request.ID, err)
		return nil, fiber.ErrInternalServerError
	}
	if len(course.ID) < 1 {
		u.Log.Warnf("Failed to find course with ID %s", request.ID)
		return nil, fiber.ErrNotFound
	}

	return course, nil
}
