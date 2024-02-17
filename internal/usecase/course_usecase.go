package usecase

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/n9mi/go-course-app/internal/entity"
	"github.com/n9mi/go-course-app/internal/helper"
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

const DEFAULT_IMAGE_URL = "https://picsum.photos/400"

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
		u.Log.Warnf("Failed to validate course : %+v", err)
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
	if err := u.CourseRepository.ScanByIDAndUserID(tx, course, request.ID, request.UserID); err != nil {
		u.Log.Warnf("Failed to find course with ID %s : %+v", request.ID, err)
		return nil, fiber.ErrInternalServerError
	}
	if len(course.ID) < 1 {
		u.Log.Warnf("Failed to find course with ID %s", request.ID)
		return nil, fiber.ErrNotFound
	}

	return course, nil
}

func (u *CourseUseCase) Create(ctx context.Context, request *model.CourseCreateRequest) (*model.CourseResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Failed to validate course : %+v", err)
		return nil, err
	}

	course := entity.Course{
		ID:          helper.GenerateRandomString(12),
		Name:        request.Name,
		Description: request.Description,
		CategoryID:  request.CategoryID,
		PriceIdr:    request.PriceIdr,
		CreatedBy:   request.CreatedBy,
		MemberCount: 0,
	}

	isDefaultImage := false
	imagePublicId := "COURSE_IMG_" + course.ID
	if request.Image != nil {
		resp, err := u.Cloudinary.Upload.Upload(ctx, request.Image, uploader.UploadParams{
			PublicID: imagePublicId,
		})
		if err != nil {
			u.Log.Warnf("Failed to upload course image : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
		course.BannerLink = resp.SecureURL
	} else {
		isDefaultImage = true
		course.BannerLink = DEFAULT_IMAGE_URL
	}

	if err := u.CourseRepository.Repository.Save(tx, &course); err != nil {
		u.Log.Warnf("Failed to insert into 'courses' : %+v", err)
		if !isDefaultImage {
			_, err := u.Cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: imagePublicId})
			if err != nil {
				u.Log.Warnf("Failed to deleting course image : %+v", err)
			}
		}
		return nil, fiber.ErrInternalServerError
	}

	response := new(model.CourseResponse)
	if err := u.CourseRepository.ScanByIDAndUserID(tx, response, course.ID, course.CreatedBy); err != nil {
		u.Log.Warnf("Failed to fetch from 'categories' : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed to commit from 'categories' : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return response, nil
}

func (u *CourseUseCase) Update(ctx context.Context, request *model.CourseUpdateRequest) (*model.CourseResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Failed to validate course : %+v", err)
		return nil, err
	}

	course := new(entity.Course)
	if err := u.CourseRepository.FindByIDAndUserID(tx, course, request.ID, request.UserID); err != nil {
		u.Log.Warnf("Course with ID %s created by %s is not found : %+v", request.ID, request.UserID, err)
		return nil, fiber.ErrNotFound
	}

	course.ID = request.ID
	course.Name = request.Name
	course.Description = request.Description
	course.CategoryID = request.CategoryID
	course.PriceIdr = request.PriceIdr

	imagePublicId := "COURSE_IMG_" + course.ID
	if request.Image != nil {
		resp, err := u.Cloudinary.Upload.Upload(ctx, request.Image, uploader.UploadParams{
			PublicID: imagePublicId,
		})
		if err != nil {
			u.Log.Warnf("Failed to upload course image : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
		course.BannerLink = resp.SecureURL
	} else if request.IsRemoveImage {
		if course.BannerLink != DEFAULT_IMAGE_URL {
			_, err := u.Cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: imagePublicId})
			if err != nil {
				u.Log.Warnf("Failed to delete course image : %+v", err)
				return nil, fiber.ErrInternalServerError
			}
			course.BannerLink = DEFAULT_IMAGE_URL
		}
	}

	if err := u.CourseRepository.Repository.Updates(tx, course); err != nil {
		u.Log.Warnf("Failed to update into 'courses' : %+v", err)
		if request.Image != nil {
			_, err := u.Cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: imagePublicId})
			if err != nil {
				u.Log.Warnf("Failed to delete course image : %+v", err)
			}
		}
		return nil, fiber.ErrInternalServerError
	}

	response := new(model.CourseResponse)
	if err := u.CourseRepository.ScanByIDAndUserID(tx, response, course.ID, course.CreatedBy); err != nil {
		u.Log.Warnf("Failed to fetch from 'courses' : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed to commit into 'courses' : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return response, nil
}

func (u *CourseUseCase) Delete(ctx context.Context, request *model.CourseDeleteRequest) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Failed to validate course : %+v", err)
		return fiber.ErrBadRequest
	}

	course := new(entity.Course)
	if err := u.CourseRepository.FindByIDAndUserID(tx, course, request.ID, request.UserID); err != nil {
		u.Log.Warnf("Course with ID %s created by %s is not found : %+v", request.ID, request.UserID, err)
		return fiber.ErrNotFound
	}

	imagePublicId := "COURSE_IMG_" + course.ID
	if _, err := u.Cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: imagePublicId}); err != nil {
		u.Log.Warnf("Failed to delete course images' : %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := u.CourseRepository.Repository.Delete(tx, course); err != nil {
		u.Log.Warnf("Failed to delete from 'courses' : %+v", err)
		return fiber.ErrInternalServerError
	}

	return nil
}
