package controller

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sec/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type SecurityController struct{}

// Login is a function that verifies the user's credentials
// and returns a JWT token if the user is valid.
//
// Parameters:
// 		- db: database connection
// 		- user: user's credentials
//
// Returns:
// 		- token: JWT token
// 		- err: error
func (sc *SecurityController) Login(db *sql.DB, mongoDB *mongo.Client) http.HandlerFunc {
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
			log.Default().Panic(err)
			return
		}

		// update last login
		accountService := AccountService{}
		_, err = accountService.UpdateLastLogin(mongoDB, user.Email)
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

func (sc *SecurityController) CreateUser(db *sql.DB, mongoDB *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		// get user's credentials
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		authUserService := AuthUserService{}
		// check if user already exist
		if authUserService.IsUserExist(db, user.Email) {
			json.NewEncoder(w).Encode(models.Error{ErrorMessage: "User already exists"})
			return
		}

		// create user's password record in postgres
		err = authUserService.CreateUser(db, user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Default().Panic(err)
			return
		}

		// create new account in mongoDB
		accountService := AccountService{}
		newAccount := models.Account{Email: user.Email, LastLogin: time.Now(), CrtDt: time.Now()}
		_, err = accountService.CreateAccount(mongoDB, newAccount)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Default().Panic(err)
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
			log.Default().Panic(err)
			return
		}
		json.NewEncoder(w).Encode(tokenString)

	}
}
