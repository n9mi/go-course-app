package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/n9mi/go-course-app/internal/model"
	"github.com/spf13/viper"
)

func NewFiber(viperConfig *viper.Viper) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      viperConfig.GetString("APP_NAME"),
		ErrorHandler: customErrorHandler(),
	})

	return app
}

func customErrorHandler() func(*fiber.Ctx, error) error {
	return func(c *fiber.Ctx, err error) error {
		var response model.WebResponse[any]

		if err != nil {
			// Check error type
			if errConv, ok := err.(validator.ValidationErrors); ok {
				response.Code = fiber.StatusBadRequest

				for _, errItem := range errConv {
					switch errItem.Tag() {
					case "required":
						response.Messages = append(response.Messages,
							fmt.Sprintf("%s is required", errItem.Field()))
					case "min":
						response.Messages = append(response.Messages,
							fmt.Sprintf("%s is should be more than %s character", errItem.Field(), errItem.Param()))
					case "max":
						response.Messages = append(response.Messages,
							fmt.Sprintf("%s is should be less than %s character", errItem.Field(), errItem.Param()))
					case "email":
						response.Messages = append(response.Messages,
							fmt.Sprintf("%s should be a valid email", errItem.Field()))
					}
				}
			} else {
				// Set default as interanl server error
				response = model.WebResponse[any]{Code: fiber.StatusInternalServerError, Messages: []string{err.Error()}}
			}

			return c.Status(response.Code).JSON(response)
		}

		return nil
	}
}
