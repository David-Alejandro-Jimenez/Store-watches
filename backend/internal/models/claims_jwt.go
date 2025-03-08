package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserName string `json:"userName"`
	jwt.RegisteredClaims
}