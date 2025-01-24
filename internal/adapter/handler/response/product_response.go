package response

type ProductResponse struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Price         int64  `json:"price"`
	StockQuantity int64  `json:"stockQuantity"`
	Category      string `json:"category"`
	CreatedAt     string `json:"createdAt"`
}
