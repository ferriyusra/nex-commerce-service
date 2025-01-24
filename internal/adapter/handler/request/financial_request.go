package request

type DepositRequest struct {
	Amount float64 `json:"amount"`
}

type WithdrawRequest struct {
	Amount float64 `json:"amount"`
}
