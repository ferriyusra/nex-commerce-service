package service

import (
	"context"
	"nex-commerce-service/internal/adapter/repository"
	"nex-commerce-service/internal/core/domain/entity"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type TransactionService interface {
	CompletePurchase(ctx context.Context, userID int64) error
}

type transactionService struct {
	transactionRepository repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{
		transactionRepository: repo,
	}
}

func (ts *transactionService) CompletePurchase(ctx context.Context, userID int64) error {
	// Ambil data cart items dari repository
	cartItems, err := ts.transactionRepository.GetCartByUserID(ctx, userID)
	if err != nil {
		code := "[SERVICE] CompletePurchase - 1"
		log.Errorw(code, err)
		return err
	}

	// Hitung totalAmount dan buat orderItems
	var totalAmount float64
	var orderItems []entity.OrderItemEntity

	for _, item := range cartItems {
		// Menghitung totalAmount
		totalAmount += float64(item.Quantity) * item.Price
		orderItems = append(orderItems, entity.OrderItemEntity{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	// Buat order baru
	order := &entity.OrderEntity{
		CustomerID: int64(userID),
		OrderDate:  time.Now(),
		Amount:     totalAmount,
	}

	if err := ts.transactionRepository.CreateOrder(ctx, order); err != nil {
		code := "[SERVICE] CompletePurchase - 2"
		log.Errorw(code, err)
		return err
	}

	// Set OrderID pada orderItems
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
	}

	// Buat order items
	if err := ts.transactionRepository.CreateOrderItems(ctx, orderItems); err != nil {
		code := "[SERVICE] CompletePurchase - 2"
		log.Errorw(code, err)
		return err
	}

	// Buat transaksi baru
	transaction := &entity.TransactionEntity{
		AccountID:       int64(userID),
		Type:            "SUCCESS",
		Amount:          totalAmount,
		TransactionDate: time.Now(),
	}

	if err := ts.transactionRepository.CreateTransaction(ctx, transaction); err != nil {
		code := "[SERVICE] CompletePurchase - 2"
		log.Errorw(code, err)
		return err
	}

	// Hapus item di keranjang setelah checkout
	if err := ts.transactionRepository.ClearCartItems(ctx, userID); err != nil {
		code := "[SERVICE] CompletePurchase - 2"
		log.Errorw(code, err)
		return err
	}

	return nil
}
