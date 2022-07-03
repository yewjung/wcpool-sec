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
	"go.mongodb.org/mongo-driver/mongo"
)

var db *sql.DB
var mongoDB *mongo.Client

func main() {
	db = driver.ConnectPostgresDB()
	mongoDB = driver.ConnectMongoDB()
	router := mux.NewRouter()

	secController := controller.SecurityController{}
	router.HandleFunc("/signin", secController.Login(db, mongoDB)).Methods("POST")
	router.HandleFunc("/signup", secController.CreateUser(db, mongoDB)).Methods("POST")
	router.HandleFunc("/refresh", secController.RefreshToken(db)).Methods("POST")

	fmt.Println("Server is running at port 8081")

	log.Fatal(
		http.ListenAndServe(":8081",
			handlers.CORS(
				handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
				handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
				handlers.AllowedOrigins([]string{"*"}),
			)(router),
		),
	)

}
