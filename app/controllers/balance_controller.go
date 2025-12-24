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

	timestamp := ctx.QueryInt("timestamp", 0)
	if timestamp == 0 {
		return utils.JsonError(ctx, errors.New("timestamp is required"), "E_TIMESTAMP_REQUIRED")
	}

	req := dto.BalanceAtTimeRequest{Timestamp: int64(timestamp)}

	return c.balanceService.GetAtTime(ctx, req, userID)
}