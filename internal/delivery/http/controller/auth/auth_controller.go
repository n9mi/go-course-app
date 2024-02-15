package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/n9mi/go-course-app/internal/model"
	"github.com/n9mi/go-course-app/internal/usecase"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	Log         *logrus.Logger
	AuthUseCase *usecase.AuthUseCase
}

func NewAuthController(authUseCase *usecase.AuthUseCase, log *logrus.Logger) *AuthController {

	return &AuthController{
		Log:         log,
		AuthUseCase: authUseCase,
	}
}

func (ct *AuthController) Register(c *fiber.Ctx) error {
	request := new(model.RegisterRequest)
	if err := c.BodyParser(request); err != nil {
		ct.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := ct.AuthUseCase.Register(c.UserContext(), request); err != nil {
		ct.Log.Warnf("Failed to register user : %+v", err)
		return err
	}

	response := model.DataResponse[any]{Data: nil}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (ct *AuthController) Login(c *fiber.Ctx) error {
	request := new(model.LoginRequest)
	if err := c.BodyParser(request); err != nil {
		ct.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	authData, err := ct.AuthUseCase.Login(c.UserContext(), request)
	if err != nil {
		ct.Log.Warnf("Failed to authenticate user : %+v", err)
		return err
	}

	// Store refresh token in HTTP only cookie
	cookie := new(fiber.Cookie)
	cookie.Name = authData.RefreshTokenName
	cookie.Value = authData.RefreshToken
	cookie.Expires = time.Unix(authData.RefreshExpAt, 0)
	cookie.HTTPOnly = true
	c.Cookie(cookie)

	response := model.DataResponse[model.TokenResponse]{Data: *authData}
	return c.Status(fiber.StatusOK).JSON(response)
}
