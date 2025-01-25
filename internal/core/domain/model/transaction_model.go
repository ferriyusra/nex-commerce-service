package model

import "time"

type Order struct {
	ID         int64     `gorm:"id"`
	CustomerID int64     `gorm:"customer_id"`
	ProductID  int64     `gorm:"product_id"`
	OrderDate  time.Time `gorm:"order_date"`
	Amount     float64   `gorm:"amount"`
}

type OrderItem struct {
	ID        int64   `gorm:"id"`
	OrderID   int64   `gorm:"order_id"`
	ProductID int64   `gorm:"product_id"`
	Quantity  int64   `gorm:"quantity"`
	Price     float64 `gorm:"price"`
}

type Transaction struct {
	ID              int64     `gorm:"id"`
	AccountID       int64     `gorm:"account_id"`
	Type            string    `gorm:"type"`
	Amount          float64   `gorm:"amount"`
	TransactionDate time.Time `gorm:"transaction_date"`
}
