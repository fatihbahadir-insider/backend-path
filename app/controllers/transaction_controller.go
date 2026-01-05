package controllers

import (
	"backend-path/app/dto"
	"backend-path/app/services"
	"backend-path/utils"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TransactionController struct {
	transactionService services.ITransactionService
}

func NewTransactionController() *TransactionController {
	return &TransactionController{
		transactionService: services.NewTransactionService(),
	}
}

func (c *TransactionController) Credit(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("user_auth").(string)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.JsonErrorUnauthorized(ctx, errors.New("invalid user id"))
	}

	var req dto.CreditRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.JsonError(ctx, err, "E_INVALID_REQUEST");
	}

	return c.transactionService.Credit(ctx, req, userID)
}

func (c *TransactionController) Debit(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("user_auth").(string)
	
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.JsonErrorUnauthorized(ctx, errors.New("invalid user id"))
	}

	var req dto.DebitRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.JsonError(ctx, err, "E_INVALID_REQUEST");
	}

	return c.transactionService.Debit(ctx, req, userID)
}

func (c *TransactionController) Transfer(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("user_auth").(string)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.JsonErrorUnauthorized(ctx, errors.New("invalid user id"))
	}

	var req dto.TransferRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.JsonError(ctx, err, "E_INVALID_REQUEST");
	}


	return c.transactionService.Transfer(ctx, req, userID)
}

func (c *TransactionController) GetByID(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("user_auth").(string)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.JsonErrorUnauthorized(ctx, errors.New("invalid user id"))
	}

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.JsonError(ctx, errors.New("invalid transaction id"), "E_INVALID_ID")
	}
	return c.transactionService.GetByID(ctx, id, userID)
}

func (c *TransactionController) GetHistory(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("user_auth").(string)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.JsonErrorUnauthorized(ctx, errors.New("invalid user id"))
	}
	return c.transactionService.GetHistory(ctx, userID)
}

func (c *TransactionController) GetStats(ctx *fiber.Ctx) error {
	return c.transactionService.GetStats(ctx)
}