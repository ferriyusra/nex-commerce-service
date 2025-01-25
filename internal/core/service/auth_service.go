package service

import (
	"context"
	"errors"
	"nex-commerce-service/config"
	"nex-commerce-service/internal/adapter/repository"
	"nex-commerce-service/internal/core/domain/entity"
	"nex-commerce-service/lib/auth"
	"nex-commerce-service/lib/conv"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/shopspring/decimal"
)

type AuthService interface {
	RegisterSeller(ctx context.Context, req entity.RegisterRequest) (*entity.UserEntity, error)
	RegisterCustomer(ctx context.Context, req entity.RegisterRequest) (*entity.UserEntity, error)
	GetUserByEmail(ctx context.Context, req entity.LoginRequest) (*entity.AccessToken, error)
}

type authService struct {
	authRepository    repository.AuthRepository
	accountRepository repository.AccountRepository
	cfg               *config.Config
	jwtToken          auth.Jwt
}

func NewAuthService(
	authRepository repository.AuthRepository,
	accountRepository repository.AccountRepository,
	cfg *config.Config, jwtToken auth.Jwt) AuthService {
	return &authService{
		authRepository:    authRepository,
		accountRepository: accountRepository,
		cfg:               cfg,
		jwtToken:          jwtToken,
	}
}

func (a *authService) RegisterSeller(ctx context.Context, req entity.RegisterRequest) (*entity.UserEntity, error) {

	password, err := conv.HashPassword(req.Password)
	if err != nil {
		code := "[SERVICE] RegisterSeller - 1"
		log.Errorw(code, err)
		return nil, err
	}

	req.Password = password

	result, err := a.authRepository.RegisterCustomer(ctx, req)
	if err != nil {
		code := "[SERVICE] RegisterSeller - 2"
		log.Errorw(code, err)
		return result, nil
	}

	err = a.accountRepository.CreateAccountWallet(ctx, entity.AccountEntity{
		UserID:  result.ID,
		Balance: decimal.NewFromFloat(0.0),
	})
	if err != nil {
		code := "[SERVICE] RegisterSeller - 3"
		log.Errorw(code, err)
		return nil, err
	}

	result = &entity.UserEntity{
		ID:       result.ID,
		Username: result.Username,
		Email:    result.Email,
		Role:     result.Role,
	}

	return result, nil
}

func (a *authService) RegisterCustomer(ctx context.Context, req entity.RegisterRequest) (*entity.UserEntity, error) {

	password, err := conv.HashPassword(req.Password)
	if err != nil {
		code := "[SERVICE] RegisterCustomer - 1"
		log.Errorw(code, err)
		return nil, err
	}

	req.Password = password

	result, err := a.authRepository.RegisterCustomer(ctx, req)
	if err != nil {
		code := "[SERVICE] RegisterCustomer - 2"
		log.Errorw(code, err)
		return result, nil
	}

	err = a.accountRepository.CreateAccountWallet(ctx, entity.AccountEntity{
		UserID:  result.ID,
		Balance: decimal.NewFromFloat(0.0),
	})
	if err != nil {
		code := "[SERVICE] RegisterCustomer - 3"
		log.Errorw(code, err)
		return nil, err
	}

	result = &entity.UserEntity{
		ID:       result.ID,
		Username: result.Username,
		Email:    result.Email,
		Role:     result.Role,
	}

	return result, nil
}

func (a *authService) GetUserByEmail(ctx context.Context, req entity.LoginRequest) (*entity.AccessToken, error) {
	result, err := a.authRepository.GetUserByEmail(ctx, req)
	if err != nil {
		code := "[SERVICE] GetUserByEmail - 1"
		log.Errorw(code, err)
		return nil, err
	}

	if checkPass := conv.CheckPasswordHash(req.Password, result.Password); !checkPass {
		code := "[SERVICE] GetUserByEmail - 2"
		err := errors.New("invalid Password")
		log.Errorw(code, err)
		return nil, err
	}

	jwtData := entity.JwtData{
		UserID:   float64(result.ID),
		RoleName: result.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
			ID:        string(result.ID),
		},
	}

	accessToken, expiresAt, err := a.jwtToken.GenerateToken(&jwtData)
	if err != nil {
		code := "[SERVICE] GetUserByEmail - 3"
		log.Errorw(code, err)
		return nil, err
	}

	response := entity.AccessToken{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}

	return &response, nil
}
