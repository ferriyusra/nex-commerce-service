package service

import (
	"context"
	"nex-commerce-service/internal/adapter/repository"
	"nex-commerce-service/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
)

type AccountService interface {
	CreateAccountWallet(ctx context.Context, req entity.AccountEntity) error
}

type accountService struct {
	accountRepository repository.AccountRepository
}

func (ac *accountService) CreateAccountWallet(ctx context.Context, req entity.AccountEntity) error {

	err := ac.accountRepository.CreateAccountWallet(ctx, req)
	if err != nil {
		code := "[SERVICE] CreateAccountWallet - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}
