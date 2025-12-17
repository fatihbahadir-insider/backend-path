package controllers

import (
	"backend-path/app/dto"
	"backend-path/app/services"
	"backend-path/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct{}

func (c *AuthController) authService() *services.AuthService {
	return new(services.AuthService)
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	utils.Logger.Info("AUTH LOGIN")

	req := new(dto.LoginRequest)
	if err := ctx.BodyParser(req); err != nil {
		return utils.JsonErrorValidation(ctx, err)
	}

	return c.authService().Authenticate(ctx, *req)
}


func (c *AuthController) Register(ctx *fiber.Ctx) error {
	utils.Logger.Info("AUTH REGISTER")

	req := new(dto.RegisterRequest)
	if err := ctx.BodyParser(req); err != nil {
		return utils.JsonErrorValidation(ctx, err)
	}

	return c.authService().Register(ctx, *req)
}