package models

import (
	"database/sql"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthUser struct {
	Email        string `json:"email"`
	PasswordHash string `json:"passwordhash"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func (user *User) VerifyUser(db *sql.DB) error {
	var authUser AuthUser

	err := db.QueryRow("SELECT username, passwordhash FROM users WHERE username=$1", user.Email).Scan(&authUser.Email, &authUser.PasswordHash)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(authUser.PasswordHash),
		[]byte(user.Password))

	if err != nil {
		return err
	}

	return nil
}

func (user *User) GenerateToken() (string, error) {
	claims := Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// get secret from os environment variable
	return token.SignedString([]byte(os.Getenv("SECRET")))
}

func (user *User) CreateUser(db *sql.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users(username, passwordhash) VALUES($1, $2)", user.Email, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}
