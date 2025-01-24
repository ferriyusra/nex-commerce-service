package entity

import "github.com/shopspring/decimal"

type AccountEntity struct {
	AccountID int64
	UserID    int64
	Balance   decimal.Decimal
}
