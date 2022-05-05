package controller

import (
	"database/sql"
	"errors"
	"os"
	"sec/constants"
	"sec/models"
	authUserRepo "sec/repository/authUser"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthUserService struct{}

// create user
func (authUserService *AuthUserService) CreateUser(db *sql.DB, user models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	authUser := models.AuthUser{
		Email:        user.Email,
		PasswordHash: string(hashedPassword),
	}

	authRepo := authUserRepo.AuthUserRepo{}
	err = authRepo.CreateUser(db, &authUser)
	if err != nil {
		return err
	}

	return nil
}

// generate token
func (authUserService *AuthUserService) GenerateToken(user models.User) (string, error) {
	claims := models.Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// get secret from os environment variable
	return token.SignedString([]byte(os.Getenv(constants.JWT_SECRET)))
}

// verify user
func (authUserService *AuthUserService) VerifyUser(db *sql.DB, user models.User) error {
	var authUser models.AuthUser

	authRepo := authUserRepo.AuthUserRepo{}
	if !authRepo.UserExist(db, user.Email) {
		// server error
		return errors.New("server error")
	}
	err := bcrypt.CompareHashAndPassword(
		[]byte(authUser.PasswordHash),
		[]byte(user.Password))

	if err != nil {
		return err
	}

	return nil
}

// refresh token
func (authUserService *AuthUserService) RefreshToken(db *sql.DB, tokenString string) (string, error) {
	claims := models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		// get secret from os environment variable
		return []byte(os.Getenv(constants.JWT_SECRET)), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}
	// if token is still valid before the expiry date, generate a new token
	if claims.ExpiresAt > time.Now().Unix() {
		return "", errors.New("token is still valid before the expiry date")
	}
	claims.ExpiresAt = time.Now().Add(time.Minute * 10).Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(os.Getenv(constants.JWT_SECRET)))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
