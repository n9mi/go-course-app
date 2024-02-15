package admin

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/n9mi/go-course-app/internal/model"
	"github.com/n9mi/go-course-app/internal/usecase"
	"github.com/sirupsen/logrus"
)

type CourseController struct {
	Log           *logrus.Logger
	CourseUseCase *usecase.CourseUseCase
}

func NewCourseController(courseUseCase *usecase.CourseUseCase, log *logrus.Logger) *CourseController {
	return &CourseController{
		Log:           log,
		CourseUseCase: courseUseCase,
	}
}

func (ct *CourseController) GetAll(c *fiber.Ctx) error {
	authData, ok := c.Locals("AuthData").(model.UserAuthData)
	if !ok {
		ct.Log.Warnf("Failed to get user auth data")
		return fiber.ErrForbidden
	}

	request := new(model.CourseListRequest)
	request.Page, _ = strconv.Atoi(c.Query("page"))
	request.PageSize, _ = strconv.Atoi(c.Query("pageSize"))
	request.UserID = authData.ID

	fmt.Println(authData)

	courses, err := ct.CourseUseCase.List(c.UserContext(), request)
	if err != nil {
		ct.Log.Warnf("Failed to get courses : %+v", err)
		return err
	}

	response := model.DataResponse[[]model.CourseListResponse]{
		Data: courses,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (ct *CourseController) GetByID(c *fiber.Ctx) error {
	authData, ok := c.Locals("AuthData").(model.UserAuthData)
	if !ok {
		ct.Log.Warnf("Failed to get user auth data")
		return fiber.ErrForbidden
	}

	request := new(model.CourseGetRequest)
	request.ID = c.Params("id")
	request.UserID = authData.ID

	course, err := ct.CourseUseCase.FindByID(c.UserContext(), request)
	if err != nil {
		ct.Log.Warnf("Failed to get course with ID %s : %+v", request.ID, err)
		return err
	}

	response := model.DataResponse[model.CourseResponse]{
		Data: *course,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
