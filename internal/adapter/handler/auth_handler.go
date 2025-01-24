package handler

import (
	"errors"
	"nex-commerce-service/internal/adapter/handler/request"
	"nex-commerce-service/internal/adapter/handler/response"
	"nex-commerce-service/internal/core/domain/entity"
	"nex-commerce-service/internal/core/service"
	validatorLib "nex-commerce-service/lib/validator"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var err error
var code string
var errorResponse response.ErrorResponseDefault
var validate = validator.New()

type AuthHandler interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
}

type authHandler struct {
	authService service.AuthService
}

// Login implements AuthHandler.
func (a *authHandler) Login(c *fiber.Ctx) error {
	req := request.LoginRequest{}
	res := response.SuccessAuthResponse{}

	if err = c.BodyParser(&req); err != nil {
		code = "[HANDLER] Login - 1"
		log.Errorw(code, err)

		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[HANDLER] Login - 2"
		log.Errorw(code, err)

		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	reqLogin := entity.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := a.authService.GetUserByEmail(c.Context(), reqLogin)
	if err != nil {
		code = "[HANDLER] Login - 3"
		log.Errorw(code, err)

		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = err.Error()

		if err.Error() == "invalid password" {
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	res.Meta.Status = true
	res.Meta.Message = "Login Successful"
	res.AccessToken = result.AccessToken
	res.ExpiresAt = result.ExpiresAt

	return c.JSON(res)
}

func (a *authHandler) Register(c *fiber.Ctx) error {
	var req request.RegisterRequest

	if err = c.BodyParser((&req)); err != nil {
		code = "[HANDLER] Register - 1"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = "Invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)

	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[HANDLER] Register - 2"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	if req.ConfirmPassword != req.Password {
		code := "[HANDLER] Register - 3"
		err = errors.New("password confirmation does not match")
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	reqEntity := entity.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Role:     "customer",
		Password: req.Password,
	}

	result, err := a.authService.Register(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] Register - 4"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	res := response.UserResponse{
		ID:       result.ID,
		Username: result.Username,
		Email:    result.Email,
		Role:     result.Role,
	}

	defaultSuccessResponse.Pagination = nil
	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Register Successful"
	defaultSuccessResponse.Data = res

	return c.JSON(defaultSuccessResponse)
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{authService: authService}
}
