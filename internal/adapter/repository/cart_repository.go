package repository

import (
	"context"
	"errors"
	"nex-commerce-service/internal/core/domain/entity"
	"nex-commerce-service/internal/core/domain/model"

	"gorm.io/gorm"
)

type CartRepository interface {
	GetCartByUserID(ctx context.Context, userID int64) ([]*entity.CartItemWithDetailsEntity, error)
	AddToCart(ctx context.Context, req entity.AddToCartEntity) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (c *cartRepository) GetCartByUserID(ctx context.Context, userID int64) ([]*entity.CartItemWithDetailsEntity, error) {
	var cartItems []*model.CartItemWithDetails

	query := `
        SELECT 
            ci.id AS cart_item_id,
            ci.product_id,
            p.name AS product_name,
            ci.quantity,
            p.price,
            p.user_id AS seller_id,
            u.username AS seller_name
        FROM cart_items ci
        JOIN cart c ON ci.cart_id = c.id
        JOIN product p ON ci.product_id = p.id
        JOIN users u ON p.user_id = u.id
        WHERE c.user_id = ?
    `

	err := c.db.Raw(query, userID).Scan(&cartItems).Error
	if err != nil {
		return nil, err
	}

	var result []*entity.CartItemWithDetailsEntity
	for _, item := range cartItems {
		result = append(result, &entity.CartItemWithDetailsEntity{
			CartItemID:  item.CartItemID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price,
			SellerID:    item.SellerID,
			SellerName:  item.SellerName,
		})
	}

	return result, nil
}

func (c *cartRepository) AddToCart(ctx context.Context, req entity.AddToCartEntity) error {
	return c.db.Transaction(func(tx *gorm.DB) error {
		// Find or create cart for user
		var cart model.Cart
		if err := tx.Table("cart").FirstOrCreate(&cart, model.Cart{UserID: req.UserID}).Error; err != nil {
			return err
		}

		// Check product exists
		var product model.Product
		if err := tx.Table("product").First(&product, "id = ?", req.ProductID).Error; err != nil {
			return err
		}

		// Validate quantity request to product stock data
		if product.StockQuantity < req.Quantity {
			return errors.New("insufficient product stock")
		}

		// Check if the product is already in the cart
		var cartItem model.CartItem
		if err := tx.Table("cart_items").Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&cartItem).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// If the product is not in the cart, create a new cart item
				cartItem = model.CartItem{
					CartID:    cart.ID,
					ProductID: req.ProductID,
					Quantity:  req.Quantity,
				}
				if err := tx.Table("cart_items").Create(&cartItem).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			// If the product is already in the cart, update the quantity
			cartItem.Quantity += req.Quantity
			if err := tx.Table("cart_items").Save(&cartItem).Error; err != nil {
				return err
			}
		}

		// Update product stock
		product.StockQuantity -= req.Quantity
		return tx.Table("product").Save(&product).Error
	})
}
