package middlewares

import (
	"backend-path/app/models"
	"backend-path/utils"
	"errors"

	"github.com/gofiber/fiber/v2"
)

func Role(allowedRoles ...models.Role) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userRole, ok := ctx.Locals("user_role").(models.Role)
		if !ok {
			return utils.JsonErrorUnauthorized(ctx, errors.New("user role not found"))
		}

		for _, role := range allowedRoles {
			if role == userRole {
				return ctx.Next()
			}
		}

		return utils.JsonErrorForbidden(ctx, errors.New("insufficient permissions"))
	}	
}