package admin

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
	authData, ok := c.Locals("AuthData").(model.UserAuthData)
	if !ok {
		ct.Log.Warn("Failed to get user auth data")
		return fiber.ErrForbidden
	}

	request := new(model.CategoryListRequest)
	request.Page, _ = strconv.Atoi(c.Query("page"))
	request.PageSize, _ = strconv.Atoi(c.Query("pageSize"))
	request.SortByPopular = strings.ToLower(c.Query("sortByPopular")) == "true"
	request.UserID = authData.ID

	categories, err := ct.CategoryUseCase.List(c.UserContext(), request)
	if err != nil {
		return err
	}

	response := model.WebResponse[[]model.CategoryResponse]{
		Data: categories,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (ct *CategoryController) Create(c *fiber.Ctx) error {
	authData, ok := c.Locals("AuthData").(model.UserAuthData)
	if !ok {
		ct.Log.Warn("Failed to get user auth data")
		return fiber.ErrForbidden
	}

	request := new(model.CategoryCreateRequest)
	if err := c.BodyParser(request); err != nil {
		ct.Log.Warnf("Failed to parsing category : %+v", err)
		return fiber.ErrBadRequest
	}
	request.UserID = authData.ID

	response, err := ct.CategoryUseCase.Create(c.UserContext(), request)
	if err != nil {
		ct.Log.Warnf("Failed to create category : %+v", err)
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (ct *CategoryController) Update(c *fiber.Ctx) error {
	categoryId := c.Params("id")

	authData, ok := c.Locals("AuthData").(model.UserAuthData)
	if !ok {
		ct.Log.Warn("Failed to get user auth data")
		return fiber.ErrForbidden
	}

	request := new(model.CategoryUpdateRequest)
	if err := c.BodyParser(request); err != nil {
		ct.Log.Warnf("Failed to parsing category : %+v", err)
		return fiber.ErrBadRequest
	}
	request.ID = categoryId
	request.UserID = authData.ID

	response, err := ct.CategoryUseCase.Update(c.UserContext(), request)
	if err != nil {
		ct.Log.Warnf("Failed to update category : %+v", err)
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (ct *CategoryController) Delete(c *fiber.Ctx) error {
	categoryId := c.Params("id")

	authData, ok := c.Locals("AuthData").(model.UserAuthData)
	if !ok {
		ct.Log.Warn("Failed to get user auth data")
		return fiber.ErrForbidden
	}

	request := new(model.CategoryDeleteRequest)
	request.ID = categoryId
	request.UserID = authData.ID

	if err := ct.CategoryUseCase.Delete(c.UserContext(), request); err != nil {
		ct.Log.Warnf("Failed to delete category : %+v", err)
		return err
	}

	return c.Status(fiber.StatusOK).JSON(model.WebResponse[any]{Data: nil})
}
