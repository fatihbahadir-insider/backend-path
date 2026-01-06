package middlewares

import (
	"backend-path/app/models"
	"backend-path/utils"
	"errors"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type JwtCustomClaims struct {
	Issuer string `json:"issuer"`
	Role models.Role `json:"role"`
	jwt.StandardClaims
}

type SkipperRoutesData struct {
	Method  string
	UrlPath string
}

func JwtMiddleware(ctx *fiber.Ctx) error {
	// skip whitelist routes
	for _, whiteList := range whiteListRoutes() {
		if ctx.Method() == whiteList.Method && ctx.Path() == whiteList.UrlPath {
			return ctx.Next()
		}
	}

	// check header token
	authorizationToken := getAuthorizationToken(ctx)
	if authorizationToken == "" {
		err := errors.New("missing Bearer token")
		return utils.JsonErrorUnauthorized(ctx, err)
	}

	// verify token
	jwtToken, err := jwt.ParseWithClaims(authorizationToken, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if jwtToken == nil {
		return utils.JsonErrorUnauthorized(ctx, errors.New("invalid token"))
	}

	claims, ok := jwtToken.Claims.(*JwtCustomClaims)
	if !ok {
		return utils.JsonErrorUnauthorized(ctx, errors.New("invalid token claims"))
	}

	ctx.Locals("user_auth", claims.Issuer)
	ctx.Locals("user_role", claims.Role)

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				ctx.Locals("token_expired", true)
				if ctx.Path() == "/api/v1/auth/refresh" {
					utils.Logger.Info("⏰ Expired token - refresh allowed")
					return ctx.Next()
				}
				return utils.JsonErrorUnauthorized(ctx, errors.New("token expired"))
			}
		}
		return utils.JsonErrorUnauthorized(ctx, err)
	}


	utils.Logger.Info("✅ SET USER AUTH")
	return ctx.Next()
}

func getAuthorizationToken(ctx *fiber.Ctx) string {
	authorizationToken := string(ctx.Request().Header.Peek("Authorization"))
	return strings.Replace(authorizationToken, "Bearer ", "", 1)
}

func whiteListRoutes() []SkipperRoutesData {
	return []SkipperRoutesData{
		{"POST", "/api/v1/auth/register"},
		{"POST", "/api/v1/auth/login"},
	}
}