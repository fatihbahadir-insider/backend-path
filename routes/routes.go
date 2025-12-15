package routes

import (
	"backend-path/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	apiRoute := app.Group("/api/v1")

	testGroup := apiRoute.Group("/test")
	testController := new(controllers.TestController)
	testGroup.Get("/", testController.List)
	
}