package middleware

import (
	"nex-commerce-service/config"
	"nex-commerce-service/internal/adapter/handler/response"
	"nex-commerce-service/lib/auth"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Middleware interface {
	CheckToken() fiber.Handler
}

type Options struct {
	authJwt auth.Jwt
}

// CheckToken implements Middleware.
func (o *Options) CheckToken() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var errorResponse response.ErrorResponseDefault
		authHandler := c.Get("Authorization")
		if authHandler == "" {
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Missing Authorization Header"
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
		}

		token := strings.Split(authHandler, "Bearer ")[1]
		claims, err := o.authJwt.VerifyAccessToken(token)
		if err != nil {
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Invalid Token"
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
		}

		c.Locals("user", claims)
		return c.Next()
	}
}

func NewMiddleware(cfg *config.Config) Middleware {
	opt := new(Options)
	opt.authJwt = auth.NewJwt(cfg)

	return opt
}