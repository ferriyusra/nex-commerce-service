package model

import "time"

type Product struct {
	ID            int64      `gorm:"id"`
	UserID        int64      `gorm:"user_id"`
	Name          string     `gorm:"name"`
	Description   string     `gorm:"description"`
	Price         float64    `gorm:"price"`
	StockQuantity int64      `gorm:"stock_quantity"`
	Category      string     `gorm:"category"`
	CreatedAt     time.Time  `gorm:"created_at"`
	UpdatedAt     *time.Time `gorm:"updated_at"`
}
