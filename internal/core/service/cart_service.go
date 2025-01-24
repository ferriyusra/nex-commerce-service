package service

import (
	"context"
	"nex-commerce-service/internal/adapter/repository"
	"nex-commerce-service/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
)

type CartService interface {
	GetCartByUserID(ctx context.Context, userID int64) ([]*entity.CartItemWithDetailsEntity, error)
	AddToCart(ctx context.Context, req entity.AddToCartEntity) error
}

type cartService struct {
	cartRepository repository.CartRepository
}

func NewCartService(repo repository.CartRepository) CartService {
	return &cartService{
		cartRepository: repo,
	}
}

func (c *cartService) GetCartByUserID(ctx context.Context, userID int64) ([]*entity.CartItemWithDetailsEntity, error) {

	result, err := c.cartRepository.GetCartByUserID(ctx, userID)
	if err != nil {
		code := "[SERVICE] GetCartByUserID - 1"
		log.Errorw(code, err)

		return nil, err
	}

	return result, nil
}

// AddToCart implements CartService.
func (c *cartService) AddToCart(ctx context.Context, req entity.AddToCartEntity) error {
	err := c.cartRepository.AddToCart(ctx, req)
	if err != nil {
		code := "[SERVICE] AddToCart - 1"
		log.Errorw(code, err)
		return err
	}
	return nil
}
