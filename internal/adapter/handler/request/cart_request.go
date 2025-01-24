package request

type AddToCartRequest struct {
	ProductID int64 `json:"productId" validate:"required"`
	Quantity  int64 `json:"quantity" validate:"required"`
}
