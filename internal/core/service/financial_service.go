package service

import (
	"context"
	"errors"
	"nex-commerce-service/internal/adapter/repository"
	"nex-commerce-service/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
	"github.com/shopspring/decimal"
)

type FinancialService interface {
	Deposit(ctx context.Context, req entity.DepositEntity) error
	Withdraw(ctx context.Context, req entity.WithdrawEntity) error
	GetBalance(ctx context.Context, userID int64) (*entity.AccountEntity, error)
}

type financialService struct {
	financialRepository repository.FinancialRepository
}

func NewFinancialService(repo repository.FinancialRepository) FinancialService {
	return &financialService{
		financialRepository: repo,
	}
}

// Deposit implements FinancialService.
func (f *financialService) Deposit(ctx context.Context, req entity.DepositEntity) error {
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		code := "[SERVICE] Deposit - 1"
		log.Errorw(code, errors.New("deposit amount must be greater than zero"))
		return nil
	}

	// ambil account balance
	account, err := f.financialRepository.GetBalance(ctx, req.UserID)
	if err != nil {
		code := "[SERVICE] Deposit - 2"
		log.Errorw(code, err)
		return err
	}

	// hitung balance terbaru
	newBalance := account.Balance.Add(req.Amount)

	// Log the perhitungan balance
	log.Debugw("Calculated New Balance", "currentBalance", account.Balance, "depositAmount", req.Amount, "newBalance", newBalance)

	err = f.financialRepository.Deposit(ctx, req)
	if err != nil {
		code := "[SERVICE] Deposit - 3"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// Withdraw implements FinancialService.
func (f *financialService) Withdraw(ctx context.Context, req entity.WithdrawEntity) error {

	if req.Amount.LessThanOrEqual(decimal.Zero) {
		code := "[SERVICE] Withdraw - 1"
		log.Errorw(code, errors.New("Withdraw amount must be greater than zero"))
		return nil
	}

	account, err := f.financialRepository.GetBalance(ctx, req.UserID)
	if err != nil {
		code := "[SERVICE] Withdraw - 2"
		log.Errorw(code, err)
		return err
	}

	// Validate balance
	if account.Balance.LessThan(req.Amount) {
		code := "[SERVICE] Withdraw - 3"
		log.Errorw(code, errors.New("insufficient balance for withdrawal"))
		return nil
	}

	// Log the perhitungan balance
	newBalance := account.Balance.Sub(req.Amount)
	log.Debugw("[SERVICE] Withdraw - 4 -Calculated New Balance", "userID", req.UserID, "currentBalance", account.Balance, "newBalance", newBalance)

	err = f.financialRepository.Withdraw(ctx, req)
	if err != nil {
		code := "[SERVICE] Withdraw - 5"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// GetBalance implements FinancialService.
func (f *financialService) GetBalance(ctx context.Context, userID int64) (*entity.AccountEntity, error) {
	result, err := f.financialRepository.GetBalance(ctx, userID)
	if err != nil {
		code := "[SERVICE] GetBalance - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return result, nil
}
