package driver

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectPostgresDB() *sql.DB {
	pgurl := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"passwordDB", 5432, "user", "mysecretpassword", "user")
	db, err := sql.Open("postgres", pgurl)
	logFatal(err)

	err = db.Ping()
	logFatal(err)

	return db
}

func ConnectMongoDB() *mongo.Client {
	uri := "mongodb://root:example@userdb:27017/"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("MongoDB successfully connected and pinged.")
	return client
}
