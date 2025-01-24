package model

import "github.com/shopspring/decimal"

type Account struct {
	ID      int64           `gorm:"id"`
	UserID  int64           `gorm:"user_id"`
	Balance decimal.Decimal `gorm:"balance"`
	Version int             `gorm:"version"`
}
