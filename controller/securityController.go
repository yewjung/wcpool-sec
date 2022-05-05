package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sec/models"
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
		authUserService := AuthUserService{}
		err = authUserService.VerifyUser(db, user)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// generate JWT token
		token, err := authUserService.GenerateToken(user)
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
		authUserService := AuthUserService{}
		err = authUserService.CreateUser(db, user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// generate JWT token
		token, err := authUserService.GenerateToken(user)
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

		authUserService := AuthUserService{}
		tokenString, err := authUserService.RefreshToken(db, tokenString)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(tokenString)

	}
}
