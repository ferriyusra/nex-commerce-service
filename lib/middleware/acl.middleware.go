package middleware

import (
	"nex-commerce-service/internal/adapter/handler/response"

	"github.com/gofiber/fiber/v2"
)

func ACLMiddleware(allowedRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var errorResponse response.ErrorResponseDefault
		role := c.Locals("roleName")
		if role == nil {
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Forbidden"
			return c.Status(fiber.StatusForbidden).JSON(errorResponse)
		}

		userRole, ok := role.(string)
		if !ok || !contains(allowedRoles, userRole) {
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Forbidden"
			return c.Status(fiber.StatusForbidden).JSON(errorResponse)
		}

		return c.Next()
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
