package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	controller "sec/controller"
	"sec/driver"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var db *sql.DB

func main() {
	db = driver.ConnectDB()
	router := mux.NewRouter()

	secController := controller.SecurityController{}
	router.HandleFunc("/signin", secController.VerifyUser(db)).Methods("POST")
	router.HandleFunc("/signup", secController.CreateUser(db)).Methods("POST")
	router.HandleFunc("/refresh", secController.RefreshToken(db)).Methods("POST")

	fmt.Println("Server is running at port 8081")

	log.Fatal(
		http.ListenAndServeTLS(":8081", "cert.pem", "key.pem",
			handlers.CORS(
				handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
				handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
				handlers.AllowedOrigins([]string{"*"}),
			)(router),
		),
	)

}
