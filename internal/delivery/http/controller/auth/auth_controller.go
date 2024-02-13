package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/n9mi/go-course-app/internal/model"
	"github.com/n9mi/go-course-app/internal/usecase"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	AuthUseCase *usecase.AuthUseCase
	Log         *logrus.Logger
}

func NewAuthController(authUseCase *usecase.AuthUseCase, log *logrus.Logger) *AuthController {

	return &AuthController{
		AuthUseCase: authUseCase,
		Log:         log,
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

	return c.Status(fiber.StatusOK).JSON(model.WebResponse[any]{Data: nil})
}

func (ct *AuthController) Login(c *fiber.Ctx) error {
	request := new(model.LoginRequest)
	if err := c.BodyParser(request); err != nil {
		ct.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response, err := ct.AuthUseCase.Login(c.UserContext(), request)
	if err != nil {
		return err
	}

	// Store refresh token in HTTP only cookie
	cookie := new(fiber.Cookie)
	cookie.Name = response.RefreshTokenName
	cookie.Value = response.RefreshToken
	cookie.Expires = time.Unix(response.RefreshExpAt, 0)
	cookie.HTTPOnly = true
	c.Cookie(cookie)

	return c.Status(fiber.StatusOK).JSON(model.WebResponse[any]{Data: response})
}
