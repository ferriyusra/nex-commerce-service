package handler

import (
	"nex-commerce-service/internal/adapter/handler/request"
	"nex-commerce-service/internal/adapter/handler/response"
	"nex-commerce-service/internal/core/domain/entity"
	"nex-commerce-service/internal/core/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/shopspring/decimal"
)

type FinancialHandler interface {
	Deposit(c *fiber.Ctx) error
	Withdraw(c *fiber.Ctx) error
	GetBalance(c *fiber.Ctx) error
}

type financialHandler struct {
	financialService service.FinancialService
}

func NewFinancialHandler(financialService service.FinancialService) FinancialHandler {
	return &financialHandler{
		financialService: financialService,
	}
}

func (fh *financialHandler) Deposit(c *fiber.Ctx) error {
	var req request.DepositRequest
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] Deposit - 1"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
	}

	if err = c.BodyParser((&req)); err != nil {
		code = "[HANDLER] Deposit - 2"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = "Invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)

	}

	reqEntity := entity.DepositEntity{
		UserID: int64(claims.UserID),
		Amount: decimal.NewFromFloat(req.Amount),
	}

	err := fh.financialService.Deposit(c.Context(), reqEntity)
	if err != nil {
		code := "[HANDLER] Deposit - 3"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	result, err := fh.financialService.GetBalance(c.Context(), int64(claims.UserID))
	if err != nil {
		code := "[HANDLER] Deposit - 4"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	res := response.UserBalanceResponse{
		AccountID: result.AccountID,
		UserID:    result.UserID,
		Balance:   result.Balance.InexactFloat64(),
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success Deposit"
	defaultSuccessResponse.Data = res

	return c.Status(fiber.StatusOK).JSON(defaultSuccessResponse)
}

func (fh *financialHandler) Withdraw(c *fiber.Ctx) error {
	var req request.WithdrawRequest
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] Withdraw - 1"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
	}

	if err = c.BodyParser((&req)); err != nil {
		code = "[HANDLER] Withdraw - 2"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = "Invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)

	}

	reqEntity := entity.WithdrawEntity{
		UserID: int64(claims.UserID),
		Amount: decimal.NewFromFloat(req.Amount),
	}

	err := fh.financialService.Withdraw(c.Context(), reqEntity)
	if err != nil {
		code := "[HANDLER] Withdraw - 3"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	result, err := fh.financialService.GetBalance(c.Context(), int64(claims.UserID))
	if err != nil {
		code := "[HANDLER] Withdraw - 4"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	res := response.UserBalanceResponse{
		AccountID: result.AccountID,
		UserID:    result.UserID,
		Balance:   result.Balance.InexactFloat64(),
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success Withdraw"
	defaultSuccessResponse.Data = res

	return c.Status(fiber.StatusOK).JSON(defaultSuccessResponse)
}

func (fh *financialHandler) GetBalance(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] GetBalance - 1"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
	}

	result, err := fh.financialService.GetBalance(c.Context(), int64(claims.UserID))
	if err != nil {
		code := "[HANDLER] GetBalance - 2"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	res := response.UserBalanceResponse{
		AccountID: result.AccountID,
		UserID:    result.UserID,
		Balance:   result.Balance.InexactFloat64(),
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success Get User Balance"
	defaultSuccessResponse.Data = res

	return c.Status(fiber.StatusOK).JSON(defaultSuccessResponse)
}
