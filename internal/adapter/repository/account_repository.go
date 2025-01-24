package repository

import (
	"context"
	"errors"
	"nex-commerce-service/internal/core/domain/entity"
	"nex-commerce-service/internal/core/domain/model"

	"github.com/gofiber/fiber/v2/log"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AccountRepository interface {
	CreateAccountWallet(ctx context.Context, req entity.AccountEntity) error
	// Deposit(ctx context.Context, req entity.DepositEntity) error
	// Withdraw(ctx context.Context, req entity.WithdrawEntity) error
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
		Balance: decimal.NewFromFloat(0.0),
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

// func (ac *accountRepository) Deposit(ctx context.Context, req entity.DepositEntity) error {
// 	return ac.db.Table("account").Transaction(func(tx *gorm.DB) error {
// 		var account model.Account
// 		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&account, req.AccountID).Error; err != nil {
// 			return err
// 		}

// 		amount := req.Amount

// 		if amount.LessThanOrEqual(decimal.Zero) {
// 			return errors.New("invalid deposit amount")
// 		}

// 		// Calculate new balance
// 		newBalance := account.Balance.Add(amount)

// 		// Update account with optimistic locking
// 		result := tx.Model(&account).
// 			Where("version = ?", account.Version).
// 			Updates(map[string]interface{}{
// 				"balance":    newBalance,
// 				"version":    gorm.Expr("version + 1"),
// 				"updated_at": time.Now(),
// 			})

// 		if result.Error != nil {
// 			return result.Error
// 		}

// 		// Prevent concurrent updates
// 		if result.RowsAffected == 0 {
// 			return errors.New("concurrent update conflict")
// 		}

// 		// Log transaction
// 		transaction := entity.AccountTransactionEntity{
// 			AccountID: req.AccountID,
// 			Amount:    amount,
// 			Type:      "DEPOSIT",
// 			Status:    "COMPLETED",
// 		}
// 		return tx.Create(&transaction).Error
// 	})
// }

// func (ac *accountRepository) Withdraw(ctx context.Context, req entity.WithdrawEntity) error {
// 	return ac.db.Table("account").Transaction(func(tx *gorm.DB) error {
// 		var account model.Account
// 		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&account, req.AccountID).Error; err != nil {
// 			return err
// 		}

// 		amount := req.Amount

// 		if amount.LessThanOrEqual(decimal.Zero) {
// 			return errors.New("invalid deposit amount")
// 		}

// 		// Calculate new balance
// 		newBalance := account.Balance.Add(amount)

// 		// Update account with optimistic locking
// 		result := tx.Model(&account).
// 			Where("version = ?", account.Version).
// 			Updates(map[string]interface{}{
// 				"balance":    newBalance,
// 				"version":    gorm.Expr("version + 1"),
// 				"updated_at": time.Now(),
// 			})

// 		if result.Error != nil {
// 			return result.Error
// 		}

// 		// Prevent concurrent updates
// 		if result.RowsAffected == 0 {
// 			return errors.New("concurrent update conflict")
// 		}

// 		// Log transaction
// 		transaction := entity.AccountTransactionEntity{
// 			AccountID: int64(req.AccountID),
// 			Amount:    amount,
// 			Type:      "DEPOSIT",
// 			Status:    "COMPLETED",
// 		}
// 		return tx.Create(&transaction).Error
// 	})
// }
