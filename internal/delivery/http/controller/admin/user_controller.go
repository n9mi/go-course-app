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
	request.RoleIDs = strings.ToLower(c.Query("role_ids"))
	if len(request.RoleIDs) > 0 && strings.Contains(request.RoleIDs, ",") {
		request.FilterRoleID = strings.Split(request.RoleIDs, ",")
	} else if len(request.RoleIDs) > 0 {
		request.FilterRoleID = []string{request.RoleIDs}
	} else {
		request.FilterRoleID = make([]string, 0)
	}

	users, err := ct.UserUseCase.List(c.UserContext(), request)
	if err != nil {
		ct.Log.Warnf("Failed to get users : %+v", err)
		return err
	}

	response := model.DataResponse[[]model.UserResponse]{
		Data: users,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (ct *UserController) UpdateRoles(c *fiber.Ctx) error {
	request := new(model.UserUpdateRolesRequest)
	request.UserID = c.Params("id")

	if err := c.BodyParser(request); err != nil {
		ct.Log.Warnf("Failed to parsing request : %+v", err)
		return fiber.ErrBadRequest
	}

	user, err := ct.UserUseCase.UpdateRoles(c.UserContext(), request)
	if err != nil {
		ct.Log.Warnf("Failed to update roles for user with ID %s : %+v", request.UserID, err)
		return err
	}

	response := model.DataResponse[model.UserResponse]{
		Data: *user,
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
