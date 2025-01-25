package entity

import "github.com/golang-jwt/jwt/v5"

type JwtData struct {
	UserID   float64 `json:"userId"`
	RoleName string  `json:"roleName"`
	jwt.RegisteredClaims
}
