package middlewares

import (
	"backend-path/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func RequestTracker() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		duration := time.Since(start)
		requestID := GetRequestID(c)

		userID := ""
		if userAuth := c.Locals("user_auth"); userAuth != nil {
			if id, ok := userAuth.(string); ok {
				userID = id
			}
		}

		ip := getClientIP(c)

		logFields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Int64("duration_ms", duration.Milliseconds()),
			zap.String("ip", ip),
			zap.String("user_agent", string(c.Request().Header.UserAgent())),
			zap.String("referer", string(c.Request().Header.Referer())),
		}

		if traceID, ok := c.Locals("trace_id").(string); ok && traceID != "" {
			logFields = append(logFields, zap.String("trace_id", traceID))
		}
		if spanID, ok := c.Locals("span_id").(string); ok && spanID != "" {
			logFields = append(logFields, zap.String("span_id", spanID))
		}

		if userID != "" {
			logFields = append(logFields, zap.String("user_id", userID))
		}
		
		if err != nil {
			logFields = append(logFields, zap.Error(err))
		}

		if c.Response().StatusCode() >= 500 {
			utils.Logger.Error("HTTP Request Error", logFields...)
		} else if c.Response().StatusCode() >= 400 {
			utils.Logger.Warn("HTTP Request Warning", logFields...)
		} else {
			utils.Logger.Info("HTTP Request", logFields...)
		}

		c.Set("X-Response-Time", duration.String())

		return err
	}
}

func getClientIP(c *fiber.Ctx) string {
	if ip := c.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	
	if ip := c.Get("X-Real-IP"); ip != "" {
		return ip
	}
	
	return c.IP()
}