package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/n9mi/go-course-app/internal/delivery/http/controller"
	"github.com/n9mi/go-course-app/internal/delivery/http/middleware"
)

type RouteConfig struct {
	App             *fiber.App
	MiddlewareSetup *middleware.MiddlewareSetup
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
	c.SetupAdminRoute(route)
	c.SetupUserRoute(route)
}

func (c *RouteConfig) SetupAuthRoute(route fiber.Router) {
	auth := route.Group("/auth")
	auth.Post("/register", c.ControllerSetup.AuthController.Register)
	auth.Post("/login", c.ControllerSetup.AuthController.Login)
}

func (c *RouteConfig) SetupAdminRoute(route fiber.Router) {
	admin := route.Group("/admin")
	admin.Use(c.MiddlewareSetup.AuthMiddleware)
	admin.Post("/categories", c.ControllerSetup.CategoryController.Create)
}

func (c *RouteConfig) SetupUserRoute(route fiber.Router) {
	admin := route.Group("/user")
	admin.Use(c.MiddlewareSetup.AuthMiddleware)
	admin.Get("/categories", c.ControllerSetup.CategoryController.Get)
	admin.Get("/categories/:id", c.ControllerSetup.CategoryController.Get)
}
