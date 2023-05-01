package auth

import (
	"github.com/dgrijalva/jwt-go"
)

type JwtClaims struct {
	UserID   string
	UserName string
	UserRole string
}

type tokenClaims struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	UserRole string `json:"user_role"`
	jwt.StandardClaims
}
