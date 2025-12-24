package controllers

import (
	"backend-path/app/dto"
	"backend-path/app/services"
	"backend-path/utils"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type BalanceController struct {
	balanceService services.IBalanceService
}

func NewBalanceController() *BalanceController {
	return &BalanceController{
		balanceService: services.NewBalanceService(),
	}
}

func (c *BalanceController) GetCurrent(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("user_auth").(string)
	
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.JsonErrorUnauthorized(ctx, errors.New("invalid user id"))
	}

	return c.balanceService.GetCurrent(ctx, userID)
}

func (c *BalanceController) GetHistorical(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("user_auth").(string)
	
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.JsonErrorUnauthorized(ctx, errors.New("invalid user id"))
	}

	return c.balanceService.GetHistorical(ctx, userID)
}

func (c *BalanceController) GetAtTime(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("user_auth").(string)
	
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.JsonErrorUnauthorized(ctx, errors.New("invalid user id"))
	}

	timestamp := ctx.Query("timestamp")
	if timestamp == "" {
		return utils.JsonError(ctx, nil, "E_TIMESTAMP_REQUIRED")
	}

	req := dto.BalanceAtTimeRequest{Timestamp: timestamp}

	return c.balanceService.GetAtTime(ctx, req, userID)
}