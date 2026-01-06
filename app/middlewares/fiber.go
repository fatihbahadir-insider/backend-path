package middlewares

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Setup(app *fiber.App) {
	app.Use(recover.New())

	app.Use(SetupHelmet())

	app.Use(SetupCors())

	app.Use(SetupRateLimiter())

	maxRequest, _ := strconv.Atoi(os.Getenv("APP_MAX_REQUEST"))
	app.Use(Limit(maxRequest, 5))

	app.Use(logger.New(logger.Config{
		Format:     `${time} ${locals:requestid} ${status} - ${method} ${url}` + "\n\n",
		TimeFormat: "2006/01/02 15:04:05",
	}))

	debug, _ := strconv.ParseBool(os.Getenv("APP_DEBUG"))
	if debug {
		app.Use(pprof.New())
	}
}

func Limit(maxRequest int, duration time.Duration) func(*fiber.Ctx) error {
	return limiter.New(limiter.Config{
		Max:        maxRequest,
		Expiration: duration,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   true,
				"message": "Too many requests",
			})
		},
	})
}

func MaxBodySize(sizeInMb int) func(*fiber.Ctx) error {
	size := sizeInMb * 1024 * 1024
	return func(c *fiber.Ctx) error {
		if len(c.Body()) >= size {
			return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
				"error":   true,
				"message": "Request body too large",
			})
		}
		return c.Next()
	}
}

