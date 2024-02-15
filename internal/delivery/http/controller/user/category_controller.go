package user

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/n9mi/go-course-app/internal/model"
	"github.com/n9mi/go-course-app/internal/usecase"
	"github.com/sirupsen/logrus"
)

type CategoryController struct {
	Log             *logrus.Logger
	CategoryUseCase *usecase.CategoryUseCase
}

func NewCategoryController(categoryUseCase *usecase.CategoryUseCase, log *logrus.Logger) *CategoryController {
	return &CategoryController{
		Log:             log,
		CategoryUseCase: categoryUseCase,
	}
}

func (ct *CategoryController) GetAll(c *fiber.Ctx) error {
	request := new(model.CategoryListRequest)
	request.Page, _ = strconv.Atoi(c.Query("page"))
	request.PageSize, _ = strconv.Atoi(c.Query("pageSize"))
	request.SortByPopular = strings.ToLower(c.Query("sortByPopular")) == "true"

	categories, err := ct.CategoryUseCase.List(c.UserContext(), request)
	if err != nil {
		ct.Log.Warnf("Failed to get categories : %+v", err)
		return err
	}

	response := model.DataResponse[[]model.CategoryResponse]{
		Data: categories,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
