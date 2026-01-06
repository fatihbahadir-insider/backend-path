package middlewares

import (
	"backend-path/configs"
	"backend-path/utils"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"go.uber.org/zap"
)

type RateLimitConfig struct {
	Max int
	Expiration time.Duration
	Skip func(*fiber.Ctx) bool
}

func SetupRateLimiter() fiber.Handler {
	maxRequest, _ := strconv.Atoi(os.Getenv("APP_MAX_REQUEST"))
	if maxRequest == 0 {
		maxRequest = 100
	}

	expirationMinutes, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_EXPIRATION_MINUTES"))
	if expirationMinutes == 0 {
		expirationMinutes = 1
	}

	config := limiter.Config{
		Max: maxRequest,
		Expiration: time.Duration(expirationMinutes) * time.Minute,
		Storage: configs.RedisStorage,
		KeyGenerator: func(c *fiber.Ctx) string {
			ip := c.IP()

			if forwarederFor := c.Get("X-Forwarded-For"); forwarederFor != "" {
				ips := strings.Split(forwarederFor, ",")
				if len(ips) > 0 {
					ip = strings.TrimSpace(ips[0])
				}
			}

			return  "rate_limit:" + ip
		},
		LimitReached: func(c *fiber.Ctx) error {
			utils.Logger.Warn("Rate limit exceeded for IP: ", zap.String("ip", c.IP()))
			return utils.JsonErrorRateLimit(c, errors.New("rate limit exceeded"))
		},
		SkipFailedRequests: false,
		SkipSuccessfulRequests: false,
	}

	return limiter.New(config)
}

func SetupAuthRateLimiter() fiber.Handler {
	maxRequest, _ := strconv.Atoi(os.Getenv("AUTH_RATE_LIMIT_MAX"))
	if maxRequest == 0 {
		maxRequest = 5
	}

	expirationMinutes, _ := strconv.Atoi(os.Getenv("AUTH_RATE_LIMIT_EXPIRATION_MINUTES"))
	if expirationMinutes == 0 {
		expirationMinutes = 1
	}
	

	return limiter.New(limiter.Config{
		Max:        maxRequest,
		Expiration: time.Duration(expirationMinutes) * time.Minute,
		Storage:    configs.RedisStorage,
		KeyGenerator: func(c *fiber.Ctx) string {
			ip := c.IP()
			if forwardedFor := c.Get("X-Forwarded-For"); forwardedFor != "" {
				ips := strings.Split(forwardedFor, ",")
				if len(ips) > 0 {
					ip = strings.TrimSpace(ips[0])
				}
			}
			return "auth_rate_limit:" + ip + ":" + c.Path()
		},
		LimitReached: func(c *fiber.Ctx) error {
			utils.Logger.Warn("Auth rate limit exceeded for IP: " + c.IP() + " on path: " + c.Path())
			return utils.JsonErrorRateLimit(c, errors.New("auth rate limit exceeded"))
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
	})
}

func SetupCustomRateLimiter(config RateLimitConfig) fiber.Handler {
	limiterConfig := limiter.Config{
		Max:        config.Max,
		Expiration: config.Expiration,
		Storage:    configs.RedisStorage,
		KeyGenerator: func(c *fiber.Ctx) string {
			ip := c.IP()
			if forwardedFor := c.Get("X-Forwarded-For"); forwardedFor != "" {
				ips := strings.Split(forwardedFor, ",")
				if len(ips) > 0 {
					ip = strings.TrimSpace(ips[0])
				}
			}
			return "rate_limit:" + ip
		},
		LimitReached: func(c *fiber.Ctx) error {
			utils.Logger.Warn("Rate limit exceeded for IP: ", zap.String("ip", c.IP()))
			return utils.JsonErrorRateLimit(c, errors.New("rate limit exceeded"))
		},
	}

	return limiter.New(limiterConfig)
}
