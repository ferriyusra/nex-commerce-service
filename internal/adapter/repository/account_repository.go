package repository

import (
	"context"
	"errors"
	"nex-commerce-service/internal/core/domain/entity"
	"nex-commerce-service/internal/core/domain/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type AccountRepository interface {
	CreateAccountWallet(ctx context.Context, req entity.AccountEntity) error
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (ac *accountRepository) CreateAccountWallet(ctx context.Context, req entity.AccountEntity) error {
	accountModel := model.Account{
		UserID:  req.UserID,
		Balance: 0,
	}

	err := ac.db.Table("account").Create(&accountModel)
	if err.Error != nil {
		if errors.Is(err.Error, gorm.ErrDuplicatedKey) {
			code = "[REPOSITORY] Create Account Wallet - 1"
			log.Errorw(code, err.Error)
			return err.Error
		}
		code = "[REPOSITORY] Create Account Wallet - 2"
		log.Errorw(code, err)
		return err.Error
	}

	return nil
}
