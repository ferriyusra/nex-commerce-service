package request

type ProductRequest struct {
	Name          string `json:"name" validate:"required"`
	Price         int64  `json:"price" validate:"required,gte=1"`
	StockQuantity int64  `json:"stockQuantity" validate:"required,gte=0"`
	Category      string `json:"category" validate:"required"`
}
