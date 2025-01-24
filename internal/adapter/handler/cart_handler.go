package handler

import (
	"nex-commerce-service/internal/adapter/handler/request"
	"nex-commerce-service/internal/adapter/handler/response"
	"nex-commerce-service/internal/core/domain/entity"
	"nex-commerce-service/internal/core/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type CartHandler interface {
	GetCartByUserID(c *fiber.Ctx) error
	AddToCart(c *fiber.Ctx) error
}

type cartHandler struct {
	cartService service.CartService
}

func NewCartHandler(cartService service.CartService) CartHandler {
	return &cartHandler{
		cartService: cartService,
	}
}

func (ct *cartHandler) GetCartByUserID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] GetCartByUserID - 1"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
	}

	result, err := ct.cartService.GetCartByUserID(c.Context(), int64(claims.UserID))
	if err != nil {
		code := "[HANDLER] GetCartByUserID - 3"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	resProducts := []response.CartItemWithDetailsResponnse{}
	for _, item := range result {
		resProduct := response.CartItemWithDetailsResponnse{
			CartItemID: item.CartItemID,
			Product: response.ProductInfo{
				ProductID:   item.ProductID,
				ProductName: item.ProductName,
				Price:       item.Price,
			},
			Quantity: item.Quantity,
			Seller: response.SellerInfo{
				SellerID:   item.SellerID,
				SellerName: item.SellerName,
			},
		}

		resProducts = append(resProducts, resProduct)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success List Carts"
	defaultSuccessResponse.Data = resProducts
	defaultSuccessResponse.Pagination = nil

	return c.JSON(defaultSuccessResponse)
}

func (ct *cartHandler) AddToCart(c *fiber.Ctx) error {
	var req request.AddToCartRequest
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] AddToCart - 1"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
	}

	if err = c.BodyParser((&req)); err != nil {
		code = "[HANDLER] AddToCart - 2"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = "Invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)

	}

	reqEntity := entity.AddToCartEntity{
		UserID:    int64(claims.UserID),
		ProductID: int64(req.ProductID),
		Quantity:  int64(req.Quantity),
	}

	err := ct.cartService.AddToCart(c.Context(), reqEntity)
	if err != nil {
		code := "[HANDLER] AddToCart - 3"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success Add To Cart"
	defaultSuccessResponse.Data = nil

	return c.Status(fiber.StatusOK).JSON(defaultSuccessResponse)
}
