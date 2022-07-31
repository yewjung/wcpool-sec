package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sec/authorization"
	controller "sec/controller"
	"sec/driver"
	"sec/models"
	"sync"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func main() {
	db := driver.ConnectPostgresDB()
	mongoDB := driver.ConnectMongoDB()
	redisCache := driver.ConnectRedis()
	storage := models.Storage{
		PostgresUserDB:    db,
		MongoAccountDB:    mongoDB,
		RedisAccountCache: redisCache,
	}
	router := mux.NewRouter()

	secController := controller.SecurityController{Storage: storage}
	router.HandleFunc("/signin", secController.Login()).Methods("POST")
	router.HandleFunc("/signup", secController.CreateUser()).Methods("POST")
	router.HandleFunc("/refresh", secController.RefreshToken()).Methods("POST")

	var wg sync.WaitGroup
	wg.Add(2)
	go startGrpc(&wg, storage)

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
	wg.Done()

}

func startGrpc(wg *sync.WaitGroup, storage models.Storage) {
	lis, err := net.Listen("tcp", fmt.Sprintf("security:%d", 8085))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	authorization.RegisterAuthorizationServer(grpcServer, authorization.AuthorizationServerImpl{Storage: storage})
	grpcServer.Serve(lis)
	wg.Done()

}
