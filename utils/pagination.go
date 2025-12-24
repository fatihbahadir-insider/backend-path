package utils

import (
	"backend-path/app/dto"

	"github.com/gofiber/fiber/v2"
)

func GetPagination(ctx *fiber.Ctx) dto.PaginationRequest {
	pagination := dto.PaginationRequest{
		Page:  ctx.QueryInt("page", 1),
		Limit: ctx.QueryInt("limit", 20),
	}
	pagination.SetDefaults()
	return pagination
}