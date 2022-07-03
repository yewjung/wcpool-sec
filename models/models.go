package models

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthUser struct {
	Email        string `json:"email"`
	PasswordHash string `json:"passwordhash"`
}

type Account struct {
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	LastLogin time.Time `json:"lastLogin"`
	CrtDt     time.Time `json:"crtDt"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type Error struct {
	ErrorMessage string `json:"errorMessage"`
}
