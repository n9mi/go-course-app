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

	admin.Get("/categories", c.ControllerSetup.AdminCategoryController.GetAll)
	admin.Post("/categories", c.ControllerSetup.AdminCategoryController.Create)
	admin.Put("/categories/:id", c.ControllerSetup.AdminCategoryController.Update)
	admin.Delete("/categories/:id", c.ControllerSetup.AdminCategoryController.Delete)

	admin.Get("/courses", c.ControllerSetup.AdminCourseController.GetAll)
	admin.Get("/courses/:id", c.ControllerSetup.AdminCourseController.GetByID)
	admin.Post("/courses", c.ControllerSetup.AdminCourseController.Create)
	admin.Put("/courses/:id", c.ControllerSetup.AdminCourseController.Update)
	admin.Delete("/courses/:id", c.ControllerSetup.AdminCourseController.Delete)

	admin.Get("/users", c.ControllerSetup.AdminUserController.GetAll)
	admin.Put("/users/:id", c.ControllerSetup.AdminUserController.UpdateRoles)
	admin.Delete("/users/:id", c.ControllerSetup.AdminUserController.Delete)
}

func (c *RouteConfig) SetupUserRoute(route fiber.Router) {
	user := route.Group("")
	user.Use(c.MiddlewareSetup.AuthMiddleware)

	user.Get("/categories", c.ControllerSetup.UserCategoryController.GetAll)
	user.Get("/courses", c.ControllerSetup.CourseController.GetAll)
	user.Get("/courses/:id", c.ControllerSetup.CourseController.GetByID)
}
