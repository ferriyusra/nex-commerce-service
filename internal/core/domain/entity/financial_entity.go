package entity

import "github.com/shopspring/decimal"

type DepositEntity struct {
	UserID int64
	Amount decimal.Decimal
}

type WithdrawEntity struct {
	UserID int64
	Amount decimal.Decimal
}

type AccountTransactionEntity struct {
	AccountID int64
	Amount    decimal.Decimal
	Type      string
	Status    string
}
