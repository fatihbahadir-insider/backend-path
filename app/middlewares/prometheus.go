package middlewares

import (
	"backend-path/app/metrics"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() == "/metrics" {
			return c.Next()
		}

		start := time.Now()
		metrics.HttpRequestsInFlight.Inc()
	
		err := c.Next()

		metrics.HttpRequestsInFlight.Dec()
		duration := time.Since(start).Seconds()

		status := strconv.Itoa(c.Response().StatusCode())
		method := c.Method()

		routePath := c.Route().Path
		if routePath == "" {
			routePath = "unknown"
		}

		metrics.HttpRequestsTotal.WithLabelValues(method, routePath, status).Inc()
		metrics.HttpRequestDuration.WithLabelValues(method, routePath, status).Observe(duration)

		return err
	}
}