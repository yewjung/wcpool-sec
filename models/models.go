package models

import (
	"github.com/golang-jwt/jwt"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthUser struct {
	Email        string `json:"email"`
	PasswordHash string `json:"passwordhash"`
	UserId       int    `json:"userid"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}
