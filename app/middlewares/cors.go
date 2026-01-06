package middlewares

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupCors() fiber.Handler {
	allowedOrigins := os.Getenv("CORS_ORIGINS")
	var origins []string

	if allowedOrigins == "" && os.Getenv("APP_ENV") == "dev" {
		origins = []string{"*"}
	} else {
		origins = strings.Split(allowedOrigins, ",")
		for i, origin := range origins {
			origins[i] = strings.TrimSpace(origin)
		}
	}

	allowCredentials := false
	if os.Getenv("CORS_ALLOW_CREDENTIALS") == "true" {
		allowCredentials = true
		if len(origins) == 1 && origins[0] == "*" {
			origins = []string{}
		}
	}

	return cors.New(cors.Config{
		AllowOrigins:     strings.Join(origins, ","),
		AllowMethods:     "GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		ExposeHeaders:    "Content-Length, Content-Type",
		AllowCredentials: allowCredentials,
		MaxAge:           86400,
	})
}