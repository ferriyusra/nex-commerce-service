package model

type Account struct {
	ID      int64   `gorm:"id"`
	UserID  int64   `gorm:"user_id"`
	Balance float64 `gorm:"price"`
}
