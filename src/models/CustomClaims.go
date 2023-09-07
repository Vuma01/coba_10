package models

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}
