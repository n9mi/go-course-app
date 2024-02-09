package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/n9mi/go-course-app/internal/delivery/http/controller"
)

type RouteConfig struct {
	App             *fiber.App
	ControllerSetup *controller.ControllerSetup
}

func (c *RouteConfig) Setup() {
	route := c.App.Group("/api")

	route.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	c.SetupAuthRoute(route)
}

func (c *RouteConfig) SetupAuthRoute(route fiber.Router) {
	auth := route.Group("/auth")
	auth.Post("/register", c.ControllerSetup.AuthController.Register)
	auth.Post("/login", c.ControllerSetup.AuthController.Login)
}
