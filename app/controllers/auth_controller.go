package controllers

import (
	"backend-path/app/dto"
	"backend-path/app/services"
	"backend-path/utils"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthController struct{
	authService services.IAuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	utils.Logger.Info("AUTH LOGIN")

	req := new(dto.LoginRequest)
	if err := ctx.BodyParser(req); err != nil {
		return utils.JsonErrorValidation(ctx, err)
	}

	return c.authService.Authenticate(ctx, *req)
}


func (c *AuthController) Register(ctx *fiber.Ctx) error {
	utils.Logger.Info("AUTH REGISTER")

	req := new(dto.RegisterRequest)
	if err := ctx.BodyParser(req); err != nil {
		return utils.JsonErrorValidation(ctx, err)
	}

	return c.authService.Register(ctx, *req)
}

func (c *AuthController) RefreshToken(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("user_auth").(string)
	
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.JsonErrorUnauthorized(ctx, errors.New("invalid user id"))
	}

	utils.Logger.Info("AUTH REFRESH TOKEN")
	
	return c.authService.RefreshToken(ctx, userID)
}