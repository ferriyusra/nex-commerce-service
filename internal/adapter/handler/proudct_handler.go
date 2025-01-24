package handler

import (
	"nex-commerce-service/internal/adapter/handler/response"
	"nex-commerce-service/internal/core/domain/entity"
	"nex-commerce-service/internal/core/service"
	"nex-commerce-service/lib/conv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ProductHandler interface {
	FindAll(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type productHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) ProductHandler {
	return &productHandler{
		productService: productService,
	}
}

// Create implements ProductHandler.
func (ph *productHandler) Create(c *fiber.Ctx) error {
	panic("unimplemented")
}

// Delete implements ProductHandler.
func (ph *productHandler) Delete(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetByID implements ProductHandler.
func (ph *productHandler) GetByID(c *fiber.Ctx) error {
	panic("unimplemented")
}

// Update implements ProductHandler.
func (ph *productHandler) Update(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (ph *productHandler) FindAll(c *fiber.Ctx) error {

	page := 1
	if c.Query("page") != "" {
		page, err = conv.StringToInt(c.Query("page"))
		if err != nil {
			code := "[HANDLER] Find All Products - 1"
			log.Errorw(code, err)
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Invalid page number"

			return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
		}
	}

	limit := 6
	if c.Query("limit") != "" {
		limit, err = conv.StringToInt(c.Query("limit"))
		if err != nil {
			code := "[HANDLER] Find All Products - 2"
			log.Errorw(code, err)
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Invalid limit number"

			return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
		}
	}

	orderBy := "created_at"
	if c.Query("orderBy") != "" {
		orderBy = c.Query("orderBy")
	}

	orderType := "desc"
	if c.Query("orderType") != "" {
		orderType = c.Query("orderType")
	}

	search := ""
	if c.Query("search") != "" {
		search = c.Query("search")
	}

	reqEntity := entity.QueryString{
		Limit:     limit,
		Page:      page,
		OrderBy:   orderBy,
		OrderType: orderType,
		Search:    search,
	}

	results, totalData, totalPages, err := ph.productService.FindAll(c.Context(), reqEntity)
	if err != nil {
		code := "[HANDLER] Find All Products - 3"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	respProducts := []response.ProductResponse{}
	for _, product := range results {
		respProduct := response.ProductResponse{
			ID:            product.ID,
			Name:          product.Name,
			Description:   product.Description,
			StockQuantity: product.StockQuantity,
			Price:         int64(product.Price),
			CreatedAt:     product.CreatedAt.Local().Format("02 January 2006"),
		}

		respProducts = append(respProducts, respProduct)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = respProducts
	defaultSuccessResponse.Pagination = &response.PaginationResponse{
		TotalRecords: int(totalData),
		Page:         page,
		PerPage:      limit,
		TotalPages:   int(totalPages),
	}
	return c.JSON(defaultSuccessResponse)
}
