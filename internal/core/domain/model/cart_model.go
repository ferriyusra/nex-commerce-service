package model

import "time"

type Cart struct {
	ID        int64     `gorm:"id"`
	UserID    int64     `gorm:"user_id"`
	CreatedAt time.Time `gorm:"created_at"`
}

type CartItem struct {
	ID        int64 `gorm:"id"`
	CartID    int64 `gorm:"cart_id"`
	ProductID int64 `gorm:"product_id"`
	Quantity  int64 `gorm:"quantity"`
}

type CartItemWithDetails struct {
	CartItemID  int64   `gorm:"column:cart_item_id"`
	ProductID   int64   `gorm:"column:product_id"`
	ProductName string  `gorm:"column:product_name"`
	Quantity    int64   `gorm:"column:quantity"`
	Price       float64 `gorm:"column:price"`
	SellerID    int64   `gorm:"column:seller_id"`
	SellerName  string  `gorm:"column:seller_name"`
}
