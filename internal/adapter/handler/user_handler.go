package handler

import (
	"errors"
	"nex-commerce-service/internal/adapter/handler/request"
	"nex-commerce-service/internal/adapter/handler/response"
	"nex-commerce-service/internal/core/domain/entity"
	"nex-commerce-service/internal/core/service"
	validatorLib "nex-commerce-service/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UserHandler interface {
	UpdatePassword(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
}

type userHandler struct {
	userService service.UserService
}


var defaultSuccessResponse response.DefaultSucceesResponse


// GetUserByID implements UserHandler.
func (u *userHandler) GetUserByID(c *fiber.Ctx) error {

	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] GetUserByID - 1"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
	}

	user, err := u.userService.GetUserByID(c.Context(), int64(claims.UserID))
	if err != nil {
		code := "[HANDLER] GetUserByID - 2"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	response := response.UserResponse{
		ID:    user.ID,
		Username:  user.Username,
		Email: user.Email,
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = response

	return c.JSON(defaultSuccessResponse)

}

// UpdatePassword implements UserHandler.
func (u *userHandler) UpdatePassword(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] UpdatePassword - 1"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
	}

	var req request.UpdatePasswordRequest
	if err = c.BodyParser(&req); err != nil {
		code := "[HANDLER] UpdatePassword - 2"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = "Invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	if err = validatorLib.ValidateStruct(&req); err != nil {
		code := "[HANDLER] UpdatePassword - 3"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

		if req.ConfirmPassword != req.NewPassword {
		code := "[HANDLER] UpdatePassword - 4"
		err = errors.New("password confirmation does not match")
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	err = u.userService.UpdatePassword(c.Context(), req.NewPassword, int64(claims.UserID))
	if err != nil {
		code := "[HANDLER] UpdatePassword - 5"
		log.Errorw(code, err)
		errorResponse.Meta.Status = false
		errorResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)

}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{userService: userService}
}