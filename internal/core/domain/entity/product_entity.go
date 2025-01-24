package entity

import "time"

type ProductEntity struct {
	ID            int64
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      string
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}

type QueryString struct {
	Limit     int
	Page      int
	OrderBy   string
	OrderType string
	Search    string
	// Status    string
}
