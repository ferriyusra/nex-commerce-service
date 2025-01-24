package repository

import (
	"context"
	"errors"
	"nex-commerce-service/internal/core/domain/entity"
	"nex-commerce-service/internal/core/domain/model"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type FinancialRepository interface {
	Deposit(ctx context.Context, req entity.DepositEntity) error
	Withdraw(ctx context.Context, req entity.WithdrawEntity) error
	GetBalance(ctx context.Context, userID int64) (*entity.AccountEntity, error)
}

type financialRepository struct {
	db *gorm.DB
}

func NewFinancialRepository(db *gorm.DB) FinancialRepository {
	return &financialRepository{db: db}
}

func (f *financialRepository) Deposit(ctx context.Context, req entity.DepositEntity) error {
	return f.db.Table("account").Transaction(func(tx *gorm.DB) error {
		var account model.Account
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&account, req.UserID).Error; err != nil {
			code = "[REPOSITORY] Deposit - 1"
			log.Errorw(code, err)
			return err
		}

		// Ensure the deposit amount is valid
		if req.Amount.LessThanOrEqual(decimal.Zero) {
			code := "[REPOSITORY] Deposit - 2"
			log.Errorw(code, errors.New("invalid deposit amount"))
			return errors.New("invalid deposit amount")
		}

		// Calculate the new balance in the repository
		newBalance := account.Balance.Add(req.Amount)

		// Update the account balance
		result := tx.Model(&account).
			Where("version = ?", account.Version).
			Updates(map[string]interface{}{
				"balance":    newBalance,
				"version":    gorm.Expr("version + 1"),
				"updated_at": time.Now(),
			})

		if result.Error != nil {
			code := "[REPOSITORY] Deposit - 3"
			log.Errorw(code, result.Error)
			return result.Error
		}

		if result.RowsAffected == 0 {
			code := "[REPOSITORY] Deposit - 4"
			log.Errorw(code, errors.New("concurrent update conflict"))
			return errors.New("concurrent update conflict")
		}

		transaction := entity.AccountTransactionEntity{
			AccountID: int64(account.ID),
			Amount:    req.Amount,
			Type:      "DEPOSIT",
			Status:    "COMPLETED",
		}
		return tx.Table("account_transactions").Create(&transaction).Error
	})
}

func (f *financialRepository) Withdraw(ctx context.Context, req entity.WithdrawEntity) error {
	return f.db.Table("account").Transaction(func(tx *gorm.DB) error {
		var account model.Account
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&account, req.UserID).Error; err != nil {
			code = "[REPOSITORY] Withdraw - 1"
			log.Errorw(code, err)
			return err
		}

		// Validate withdrawal amount
		if req.Amount.LessThanOrEqual(decimal.Zero) {
			code := "[REPOSITORY] Withdraw - 2"
			log.Errorw(code, errors.New("invalid withdrawal amount"))
			return errors.New("invalid withdrawal amount")
		}

		// Check if the account has sufficient balance for withdrawal
		if account.Balance.LessThan(req.Amount) {
			code := "[REPOSITORY] Withdraw - 3"
			log.Errorw(code, errors.New("insufficient"))
			return errors.New("insufficient funds")
		}

		// Subtract the withdrawal amount from the balance
		newBalance := account.Balance.Sub(req.Amount)

		// Update account balance with optimistic locking
		result := tx.Model(&account).
			Where("version = ?", account.Version).
			Updates(map[string]interface{}{
				"balance":    newBalance,
				"version":    gorm.Expr("version + 1"),
				"updated_at": time.Now(),
			})

		if result.Error != nil {
			return result.Error
		}

		// Prevent concurrent updates
		if result.RowsAffected == 0 {
			code := "[REPOSITORY] Withdraw - 4"
			log.Errorw(code, errors.New("concurrent update conflict"))
			return errors.New("concurrent update conflict")
		}

		// Create transaction record
		transaction := entity.AccountTransactionEntity{
			AccountID: int64(account.ID),
			Amount:    req.Amount,
			Type:      "WITHDRAWAL",
			Status:    "COMPLETED",
		}
		return tx.Table("account_transactions").Create(&transaction).Error
	})
}

func (f *financialRepository) GetBalance(ctx context.Context, userID int64) (*entity.AccountEntity, error) {
	var user model.Account
	err = f.db.Table("account").Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		code := "[REPOSITORY] GetBalance - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return &entity.AccountEntity{
		AccountID: user.ID,
		UserID:    user.UserID,
		Balance:   user.Balance,
	}, nil

}
