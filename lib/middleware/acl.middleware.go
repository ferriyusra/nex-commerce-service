package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func ACLMiddleware(allowedRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("user_role")
		if role == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "forbidden"})
		}

		userRole, ok := role.(string)
		if !ok || !contains(allowedRoles, userRole) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "forbidden"})
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
