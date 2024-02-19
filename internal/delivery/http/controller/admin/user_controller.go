package admin

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/n9mi/go-course-app/internal/model"
	"github.com/n9mi/go-course-app/internal/usecase"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	UserUseCase *usecase.UserUseCase
	Log         *logrus.Logger
}

func NewUserController(userUseCase *usecase.UserUseCase, log *logrus.Logger) *UserController {
	return &UserController{
		UserUseCase: userUseCase,
		Log:         log,
	}
}

func (ct *UserController) GetAll(c *fiber.Ctx) error {
	request := new(model.UserListRequest)
	request.Page, _ = strconv.Atoi(c.Query("page"))
	request.PageSize, _ = strconv.Atoi(c.Query("pageSize"))
	request.SearchName = strings.ToLower(c.Query("name"))
	request.SearchEmail = strings.ToLower(c.Query("email"))
	request.FilterRoleID = strings.Split(c.Query("role_id"), ",")

	users, err := ct.UserUseCase.List(c.UserContext(), request)
	if err != nil {
		ct.Log.Warnf("Failed to get users : %+v", err)
		return err
	}

	response := model.DataResponse[[]model.UserListResponse]{
		Data: users,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (ct *UserController) Delete(c *fiber.Ctx) error {
	request := new(model.UserDeleteRequest)
	request.ID = c.Params("id")

	if err := ct.UserUseCase.Delete(c.UserContext(), request); err != nil {
		ct.Log.Warnf("Failed to delete users : %+v", err)
		return err
	}

	response := model.DataResponse[any]{
		Data: nil,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
