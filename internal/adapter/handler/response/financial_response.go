package response

type UserBalanceResponse struct {
	AccountID int64   `json:"accountId"`
	UserID    int64   `json:"userId"`
	Balance   float64 `json:"balance"`
}
