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
	CategoryUseCase *usecase.CategoryUseCase
	Log             *logrus.Logger
}

func NewCategoryController(categoryUseCase *usecase.CategoryUseCase, log *logrus.Logger) *CategoryController {
	return &CategoryController{
		CategoryUseCase: categoryUseCase,
		Log:             log,
	}
}

func (ct *CategoryController) GetAll(c *fiber.Ctx) error {
	request := new(model.CategoryListRequest)
	request.Page, _ = strconv.Atoi(c.Query("page"))
	request.PageSize, _ = strconv.Atoi(c.Query("pageSize"))
	request.SortByPopular = strings.ToLower(c.Query("sortByPopular")) == "true"

	categories, err := ct.CategoryUseCase.List(c.UserContext(), request)
	if err != nil {
		return err
	}

	response := model.WebResponse[[]model.CategoryResponse]{
		Data: categories,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
