package models

import "github.com/dgrijalva/jwt-go"

type TokenClaims struct {
	UserID   string `json:"user_id"`
	ClientIP string `json:"client_ip"`
	jwt.StandardClaims
}
