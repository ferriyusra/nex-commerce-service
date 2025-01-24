package request

type AccountRequest struct {
	Balance float64 `json:"balance" validate:"required"`
}
