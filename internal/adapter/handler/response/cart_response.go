package response

type CartItemWithDetailsResponnse struct {
	CartItemID int64       `json:"cartItemId"`
	Product    ProductInfo `json:"product"`
	Quantity   int64       `json:"quantity"`
	Seller     SellerInfo  `json:"seller"`
}

type ProductInfo struct {
	ProductID   int64   `json:"productId"`
	ProductName string  `json:"productName"`
	Price       float64 `json:"price"`
}

type SellerInfo struct {
	SellerID   int64  `json:"sellerId"`
	SellerName string `json:"sellerName"`
}
