package handler

import (
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

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{authService: authService}
}