package handler

import (
	"nex-commerce-service/internal/core/domain/entity"
	"nex-commerce-service/internal/core/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type TransactionHandler interface {
	Checkout(c *fiber.Ctx) error
}

type transactionHandler struct {
	transactionService service.TransactionService
}

func (t *transactionHandler) Checkout(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] Checkout - 1"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
	}

	cartItems := make([]entity.CartItemEntity, 0)
	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += float64(item.Quantity) * item.Product.Price
	}

	if err := t.transactionService.CompletePurchase(c.Context(), int64(claims.UserID)); err != nil {
		code := "[HANDLER] Checkout - 2"
		log.Errorw(code, err)

		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = "Internal Server Error"
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Checkout Successful"
	defaultSuccessResponse.Data = nil
	defaultSuccessResponse.Pagination = nil

	return c.JSON(defaultSuccessResponse)

}

func NewTransactionHandler(transactionService service.TransactionService) TransactionHandler {
	return &transactionHandler{
		transactionService: transactionService,
	}
}
