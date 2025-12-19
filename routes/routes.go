package routes

import (
	"backend-path/app/controllers"
	"backend-path/app/middlewares"
	"backend-path/app/models"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	apiRoute := app.Group("/api/v1", middlewares.JwtMiddleware)

	auth := apiRoute.Group("/auth")
	authController := controllers.NewAuthController()
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/refresh", authController.RefreshToken)

	users := apiRoute.Group("/users")
	userController := controllers.NewUserController()
	users.Get("/", middlewares.Role(models.RoleAdmin, models.RoleMod), userController.GetAll)
	users.Get("/:id", middlewares.Role(models.RoleAdmin, models.RoleMod), userController.GetByID)
	users.Put("/:id", middlewares.Role(models.RoleAdmin), userController.Update)
	users.Delete("/:id", middlewares.Role(models.RoleAdmin), userController.Delete)
}