package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"sec/models"
	"time"

	"github.com/golang-jwt/jwt"
)

type SecurityController struct{}

// VerifyUser is a function that verifies the user's credentials
// and returns a JWT token if the user is valid.
//
// Parameters:
// 		- db: database connection
// 		- user: user's credentials
//
// Returns:
// 		- token: JWT token
// 		- err: error
func (sc *SecurityController) VerifyUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		// get user's credentials
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// verify user's credentials
		err = user.VerifyUser(db)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// generate JWT token
		token, err := user.GenerateToken()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// return token
		json.NewEncoder(w).Encode(token)
	}
}

func (sc *SecurityController) CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		// get user's credentials
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// create user's account
		err = user.CreateUser(db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// generate JWT token
		token, err := user.GenerateToken()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// return token
		json.NewEncoder(w).Encode(token)
	}
}

func (sc *SecurityController) RefreshToken(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// retrieve token strings from request header
		tokenString := r.Header.Get("Authorization")

		// Parse token string and check if it's valid and retrieve the claim from
		claim := models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, &claim, func(token *jwt.Token) (interface{}, error) {
			// get secret from os environment variable
			return []byte(os.Getenv("SECRET")), nil
		},
		)
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// if token is still valid before the expiry date, generate a new token
		if claim.ExpiresAt > time.Now().Unix() {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		claim.ExpiresAt = time.Now().Add(time.Minute * 10).Unix()
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
		tokenString, err = token.SignedString([]byte(os.Getenv("SECRET")))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(tokenString)

	}
}
