package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"sec/models"
	"time"
)

type SecurityController struct {
	Storage models.Storage
}

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
func (sc *SecurityController) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.UserDTO

		// get user's credentials
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// verify user's credentials
		authUserService := AuthUserService{DB: sc.Storage.PostgresUserDB}
		err = authUserService.VerifyUser(user)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Default().Panic(err)
			return
		}

		// update last login
		accountService := AccountService{MongoDB: sc.Storage.MongoAccountDB}
		_, err = accountService.UpdateLastLogin(user.Email)
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

func (sc *SecurityController) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.UserDTO

		// get user's credentials
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		authUserService := AuthUserService{DB: sc.Storage.PostgresUserDB}
		// check if user already exist
		if authUserService.IsUserExist(user.Email) {
			json.NewEncoder(w).Encode(models.Error{ErrorMessage: "User already exists"})
			return
		}

		// create user's password record in postgres
		err = authUserService.CreateUser(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Default().Panic(err)
			return
		}

		// create new account in mongoDB
		accountService := AccountService{MongoDB: sc.Storage.MongoAccountDB}
		newAccount := models.Account{Email: user.Email, LastLogin: time.Now(), CrtDt: time.Now()}
		_, err = accountService.CreateAccount(newAccount)
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

func (sc *SecurityController) RefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// retrieve token strings from request header
		tokenString := r.Header.Get("Authorization")

		authUserService := AuthUserService{DB: sc.Storage.PostgresUserDB}
		tokenString, err := authUserService.RefreshToken(tokenString)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Default().Panic(err)
			return
		}
		json.NewEncoder(w).Encode(tokenString)

	}
}
