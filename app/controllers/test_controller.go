package controllers

import (
	"backend-path/app/services"
	"backend-path/utils"

	"github.com/gofiber/fiber/v2"
)

type TestController struct {}

func (c *TestController) testService() *services.TestService {
	return new(services.TestService)
}

func (c *TestController) List(ctx *fiber.Ctx) error {
	utils.Logger.Info("TEST 123")

	return c.testService().List(ctx)
}