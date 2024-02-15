package user

import (
	"strconv"
	"strings"

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
	request := new(model.CourseListRequest)
	request.Page, _ = strconv.Atoi(c.Query("page"))
	request.PageSize, _ = strconv.Atoi(c.Query("pageSize"))
	request.CategoryID = c.Query("categoryID")
	request.IsFree = strings.ToLower(c.Query("isFree")) == "true"
	request.SortByMaximumPrice = strings.ToLower(c.Query("sortByMaximumPrice")) == "true"
	request.SortByMinimumPrice = strings.ToLower(c.Query("sortByMinimumPrice")) == "true"
	request.SearchTitle = strings.ToLower(c.Query("title"))

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
	request := new(model.CourseGetRequest)
	request.ID = c.Params("id")

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
