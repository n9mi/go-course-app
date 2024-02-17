package admin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/n9mi/go-course-app/internal/helper"
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

func (ct *CourseController) Create(c *fiber.Ctx) error {
	authData, ok := c.Locals("AuthData").(model.UserAuthData)
	if !ok {
		ct.Log.Warn("Failed to get user auth data")
		return fiber.ErrForbidden
	}

	priceIdr, _ := strconv.ParseFloat(c.FormValue("price_idr"), 64)

	request := &model.CourseCreateRequest{
		Name:        c.FormValue("name"),
		Description: c.FormValue("description"),
		CategoryID:  c.FormValue("category_id"),
		PriceIdr:    priceIdr,
		CreatedBy:   authData.ID,
	}

	imageFormFile, _ := c.FormFile("image")
	if imageFormFile != nil {
		var err error
		request.Image, err = helper.FormFileToBuffer(ct.Log, imageFormFile)
		if err != nil {
			ct.Log.Warnf("Failed to process the image : %+v", err)
			return err
		}
	}

	course, err := ct.CourseUseCase.Create(c.UserContext(), request)
	if err != nil {
		ct.Log.Warnf("Failed to create course : %+v", err)
		return err
	}

	response := model.DataResponse[model.CourseResponse]{
		Data: *course,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (ct *CourseController) Update(c *fiber.Ctx) error {
	authData, ok := c.Locals("AuthData").(model.UserAuthData)
	if !ok {
		ct.Log.Warn("Failed to get user auth data")
		return fiber.ErrForbidden
	}

	courseID := c.Params("id")
	priceIdr, _ := strconv.ParseFloat(c.FormValue("price_idr"), 64)

	request := &model.CourseUpdateRequest{
		ID:            courseID,
		Name:          c.FormValue("name"),
		Description:   c.FormValue("description"),
		CategoryID:    c.FormValue("category_id"),
		PriceIdr:      priceIdr,
		UserID:        authData.ID,
		IsRemoveImage: strings.ToLower(c.FormValue("is_remove_image")) == "true",
	}

	imageFormFile, _ := c.FormFile("image")
	if imageFormFile != nil {
		var err error
		request.Image, err = helper.FormFileToBuffer(ct.Log, imageFormFile)
		if err != nil {
			ct.Log.Warnf("Failed to process the image : %+v", err)
			return err
		}
	}

	course, err := ct.CourseUseCase.Update(c.UserContext(), request)
	if err != nil {
		ct.Log.Warnf("Failed to update course : %+v", err)
		return err
	}

	response := model.DataResponse[model.CourseResponse]{
		Data: *course,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (ct *CourseController) Delete(c *fiber.Ctx) error {
	authData, ok := c.Locals("AuthData").(model.UserAuthData)
	if !ok {
		ct.Log.Warn("Failed to get user auth data")
		return fiber.ErrForbidden
	}

	request := &model.CourseDeleteRequest{
		ID:     c.Params("id"),
		UserID: authData.ID,
	}

	if err := ct.CourseUseCase.Delete(c.UserContext(), request); err != nil {
		ct.Log.Warn("Failed to delete course : %+v", err)
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusOK).JSON(model.DataResponse[any]{Data: nil})
}
