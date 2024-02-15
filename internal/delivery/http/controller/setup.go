package controller

import (
	"github.com/n9mi/go-course-app/internal/delivery/http/controller/admin"
	"github.com/n9mi/go-course-app/internal/delivery/http/controller/auth"
	"github.com/n9mi/go-course-app/internal/delivery/http/controller/user"
	"github.com/n9mi/go-course-app/internal/usecase"
	"github.com/sirupsen/logrus"
)

type ControllerSetup struct {
	AuthController *auth.AuthController

	AdminCategoryController *admin.CategoryController
	AdminCourseController   *admin.CourseController

	UserCategoryController *user.CategoryController
	UserCourseController   *user.CourseController
}

func Setup(useCaseSetup *usecase.UseCaseSetup, log *logrus.Logger) *ControllerSetup {

	return &ControllerSetup{
		// Authentication controller
		AuthController: auth.NewAuthController(useCaseSetup.AuthUseCase, log),

		// Admin controller
		AdminCategoryController: admin.NewCategoryController(useCaseSetup.CategoryUseCase, log),
		AdminCourseController:   admin.NewCourseController(useCaseSetup.CourseUseCase, log),

		// User controller
		UserCategoryController: user.NewCategoryController(useCaseSetup.CategoryUseCase, log),
		UserCourseController:   user.NewCourseController(useCaseSetup.CourseUseCase, log),
	}
}
