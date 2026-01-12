package middlewares

import (
	"backend-path/app/metrics"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Setup(app *fiber.App) {
	metrics.Init()

	app.Use(recover.New())

	app.Use(RequestID())

	app.Use(PrometheusMiddleware())

	app.Use(SetupHelmet())

	app.Use(SetupCors())

	app.Use(SetupRateLimiter())

	app.Use(logger.New(logger.Config{
		Format:     `${time} ${locals:requestid} ${status} - ${method} ${url}` + "\n\n",
		TimeFormat: "2006/01/02 15:04:05",
		Output: os.Stdout,
	}))

	app.Use(RequestTracker())

	debug, _ := strconv.ParseBool(os.Getenv("APP_DEBUG"))
	if debug {
		app.Use(pprof.New())
	}
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

