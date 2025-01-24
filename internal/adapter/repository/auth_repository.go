package repository

import (
	"context"
	"errors"
	"nex-commerce-service/internal/core/domain/entity"
	"nex-commerce-service/internal/core/domain/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

var err error
var code string

type AuthRepository interface {
	RegisterSeller(ctx context.Context, req entity.RegisterRequest) (*entity.UserEntity, error)
	RegisterCustomer(ctx context.Context, req entity.RegisterRequest) (*entity.UserEntity, error)
	GetUserByEmail(ctx context.Context, req entity.LoginRequest) (*entity.UserEntity, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (a *authRepository) GetUserByEmail(ctx context.Context, req entity.LoginRequest) (*entity.UserEntity, error) {
	var modelUser model.User

	err = a.db.Where("email = ?", req.Email).First(&modelUser).Error
	if err != nil {
		code = "[REPOSITORY] GetUserByEmail - 1"
		log.Errorw(code, err)
		return nil, err
	}

	response := entity.UserEntity{
		ID:       modelUser.ID,
		Username: modelUser.Username,
		Role:     modelUser.Role,
		Email:    modelUser.Email,
		Password: modelUser.Password,
	}

	return &response, nil
}

func (a *authRepository) RegisterCustomer(ctx context.Context, req entity.RegisterRequest) (*entity.UserEntity, error) {
	modelUser := model.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
		Password: req.Password,
	}

	err := a.db.Create(&modelUser)
	if err.Error != nil {
		if errors.Is(err.Error, gorm.ErrDuplicatedKey) {
			code = "[REPOSITORY] RegisterCustomer - 1"
			log.Errorw(code, err.Error)
			return nil, err.Error
		}
		code = "[REPOSITORY] RegisterCustomer - 2"
		log.Errorw(code, err)
		return nil, err.Error
	}

	return &entity.UserEntity{
		ID:       modelUser.ID,
		Username: modelUser.Username,
		Email:    modelUser.Email,
		Role:     modelUser.Role,
	}, nil
}

func (a *authRepository) RegisterSeller(ctx context.Context, req entity.RegisterRequest) (*entity.UserEntity, error) {
	modelUser := model.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
		Password: req.Password,
	}

	err := a.db.Create(&modelUser)
	if err.Error != nil {
		if errors.Is(err.Error, gorm.ErrDuplicatedKey) {
			code = "[REPOSITORY] RegisterSeller - 1"
			log.Errorw(code, err.Error)
			return nil, err.Error
		}
		code = "[REPOSITORY] RegisterSeller - 2"
		log.Errorw(code, err)
		return nil, err.Error
	}

	return &entity.UserEntity{
		ID:       modelUser.ID,
		Username: modelUser.Username,
		Email:    modelUser.Email,
		Role:     modelUser.Role,
	}, nil
}
