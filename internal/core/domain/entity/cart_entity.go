package entity

type AddToCartEntity struct {
	UserID    int64
	ProductID int64
	Quantity  int64
}

type CartItemWithDetailsEntity struct {
	CartItemID  int64
	ProductID   int64
	ProductName string
	Quantity    int64
	Price       float64
	SellerID    int64
	SellerName  string
}
