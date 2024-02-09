package controller

import (
	"github.com/n9mi/go-course-app/internal/usecase"
	"github.com/sirupsen/logrus"
)

type ControllerSetup struct {
	AuthController *AuthController
}

func Setup(useCaseSetup *usecase.UseCaseSetup, log *logrus.Logger) *ControllerSetup {

	return &ControllerSetup{
		AuthController: NewAuthController(useCaseSetup.AuthUseCase, log),
	}
}
