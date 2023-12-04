package model

import (
	jwt "github.com/golang-jwt/jwt/v4"
)

type SessionClaims struct {
	jwt.StandardClaims
	UserID      int64  `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
}
