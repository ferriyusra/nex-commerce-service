package repository

import (
	"context"
	"nex-commerce-service/internal/core/domain/entity"
	"nex-commerce-service/internal/core/domain/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	GetCartByUserID(ctx context.Context, userID int64) ([]*entity.CartItemWithDetailsEntity, error)
	CreateOrder(ctx context.Context, order *entity.OrderEntity) error
	CreateOrderItems(ctx context.Context, orderItems []entity.OrderItemEntity) error
	CreateTransaction(ctx context.Context, transaction *entity.TransactionEntity) error
	ClearCartItems(ctx context.Context, userID int64) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (t *transactionRepository) GetCartByUserID(ctx context.Context, userID int64) ([]*entity.CartItemWithDetailsEntity, error) {
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

	err := t.db.Raw(query, userID).Scan(&cartItems).Error
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

func (t *transactionRepository) CreateOrder(ctx context.Context, order *entity.OrderEntity) error {
	err := t.db.Table("orders").Create(order).Error

	if err != nil {
		code = "[REPOSITORY] CreateOrder - 1"
		log.Errorw(code, err.Error())
		return err
	}

	return nil
}

func (t *transactionRepository) CreateOrderItems(ctx context.Context, orderItems []entity.OrderItemEntity) error {
	err := t.db.Table("order_items").Create(&orderItems).Error

	if err != nil {
		code = "[REPOSITORY] CreateOrderItems - 1"
		log.Errorw(code, err.Error())
		return err
	}

	return nil
}

func (t *transactionRepository) CreateTransaction(ctx context.Context, transaction *entity.TransactionEntity) error {
	err := t.db.Table("transactions").Create(transaction).Error

	if err != nil {
		code = "[REPOSITORY] CreateTransaction - 1"
		log.Errorw(code, err.Error())
		return err
	}

	return nil
}

func (t *transactionRepository) ClearCartItems(ctx context.Context, userID int64) error {
	err := t.db.Exec("DELETE FROM cart_items WHERE cart_id IN (SELECT id FROM cart WHERE user_id = ?)", userID).Error

	if err != nil {
		code = "[REPOSITORY] ClearCartItems - 1"
		log.Errorw(code, err.Error())
		return err
	}

	return nil
}
