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

	UserCategoryController *user.CategoryController
}

func Setup(useCaseSetup *usecase.UseCaseSetup, log *logrus.Logger) *ControllerSetup {

	return &ControllerSetup{
		AuthController: auth.NewAuthController(useCaseSetup.AuthUseCase, log),

		AdminCategoryController: admin.NewCategoryController(useCaseSetup.CategoryUseCase, log),

		UserCategoryController: user.NewCategoryController(useCaseSetup.CategoryUseCase, log),
	}
}
