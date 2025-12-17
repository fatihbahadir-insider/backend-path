package routes

import (
	"backend-path/app/controllers"
	"backend-path/app/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	apiRoute := app.Group("/api/v1", middlewares.JwtMiddleware)

	auth := apiRoute.Group("/auth")
	authController := new(controllers.AuthController)
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)

	testGroup := apiRoute.Group("/test")
	testController := new(controllers.TestController)
	testGroup.Get("/", testController.List)

}