package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/n9mi/go-course-app/internal/model"
)

type CategoryController struct {
}

func NewCategoryController() *CategoryController {
	return &CategoryController{}
}

func (ct *CategoryController) Create(c *fiber.Ctx) error {
	authData, ok := c.Locals("AuthData").(model.UserAuthData)
	if !ok {
		return fiber.ErrUnauthorized
	}

	return c.Status(fiber.StatusOK).JSON(model.WebResponse[model.UserAuthData]{Data: authData})
}
