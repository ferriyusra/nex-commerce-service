package entity

import "time"

type CheckoutEntity struct {
	UserID int64 `json:"user_id"`
}

type OrderEntity struct {
	ID         int64     `json:"id"`
	CustomerID int64     `json:"customerId"`
	OrderDate  time.Time `json:"orderDate"`
	Amount     float64   `json:"amount"`
}

type OrderItemEntity struct {
	ID        int64   `json:"id"`
	OrderID   int64   `json:"orderId"`
	ProductID int64   `json:"productId"`
	Quantity  int64   `json:"quantity"`
	Price     float64 `json:"price"`
}

type TransactionEntity struct {
	ID              int64     `json:"id"`
	AccountID       int64     `json:"accountId"`
	Type            string    `json:"type"`
	Amount          float64   `json:"amount"`
	TransactionDate time.Time `json:"transactionDate"`
}

type CartItemEntity struct {
	ID        int64 `json:"id"`
	CartID    int64 `json:"cartId"`
	ProductID int64 `json:"productId"`
	Quantity  int   `json:"quantity"`
	Product   ProductEntity
}

type CartEntity struct {
	ID     int64            `json:"id"`
	UserID int64            `json:"userId"`
	Items  []CartItemEntity `json:"items"`
}
