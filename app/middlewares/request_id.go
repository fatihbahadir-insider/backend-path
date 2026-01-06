package middlewares

import (
	"backend-path/constants"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get(constants.RequestIDHeader)

		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set(constants.RequestIDHeader, requestID)
		c.Locals(constants.RequestIDLocal, requestID)

		return c.Next()
	}
}

func GetRequestID(c *fiber.Ctx) string {
	if id := c.Locals(constants.RequestIDLocal); id != nil {
		if idStr, ok := id.(string); ok {
			return idStr
		}
	}
	return ""
}