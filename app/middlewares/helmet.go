package middlewares

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

func SetupHelmet() fiber.Handler {
	// Note to me: Tarayıcınız bunu aldıktan sonra, belirtilen süre boyunca (örneğin 1 yıl) o siteye her erişimde otomatik olarak HTTP yerine HTTPS kullanır. 
	hstsMaxAge := 31536000
	if os.Getenv("APP_ENV") == "production" || os.Getenv("FORCE_HTTPS") == "true" {
		if customMaxAge := os.Getenv("HSTS_MAX_AGE"); customMaxAge != "" {
			if age, err := strconv.Atoi(customMaxAge); err == nil {
				hstsMaxAge = age
			}
		}
	} else {
		hstsMaxAge = 0; // Note to me: Development ortamında olduğumuz için kapadık.
	}

	csp := os.Getenv("CSP_POLICY")
	if csp == "" {
		csp = "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:;"
	}

	return helmet.New(helmet.Config{
		XSSProtection: "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions: "DENY",
		ReferrerPolicy: "strict-origin-when-cross-origin",
		CrossOriginEmbedderPolicy: "require-corp",
		CrossOriginOpenerPolicy:   "same-origin",
		CrossOriginResourcePolicy: "cross-origin",
		HSTSMaxAge:            hstsMaxAge,
		HSTSExcludeSubdomains: false,
		HSTSPreloadEnabled:    os.Getenv("APP_ENV") == "production",
		ContentSecurityPolicy: csp,
	})
}