package controllers

import (
	"backend-path/app/dto"
	"backend-path/app/services"
	"backend-path/utils"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserController struct {
	userService services.IUserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

func (c *UserController) GetAll(ctx *fiber.Ctx) error {
	return c.userService.GetAll(ctx)
}

func (c *UserController) GetByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.JsonErrorNotFound(ctx, errors.New("invalid user id"))
	}

	return c.userService.GetByID(ctx, id)
}

func (c *UserController) Update(ctx *fiber.Ctx) error {
	var req dto.UpdateUserRequest

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.JsonErrorNotFound(ctx, errors.New("invalid user id"))
	}

	if err := ctx.BodyParser(&req); err != nil {
		return utils.JsonError(ctx, err, "E_PARSE")
	}


	return c.userService.Update(ctx, id, req)
}

func (c *UserController) Delete(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.JsonErrorNotFound(ctx, errors.New("invalid user id"))
	}

	return c.userService.Delete(ctx, id)
}

func (c *UserController) GetMe(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("user_auth").(string)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.JsonErrorUnauthorized(ctx, errors.New("invalid user id"))
	}

	return c.userService.GetByID(ctx, userID)	
}