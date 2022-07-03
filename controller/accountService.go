package controller

import (
	"context"
	"sec/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountService struct{}

func (as *AccountService) CreateAccount(mongoDB *mongo.Client, account models.Account) (*mongo.InsertOneResult, error) {
	accountCollection := mongoDB.Database("Account").Collection("AccountCollection")
	doc := bson.M{"_id": account.Email, "Username": account.Username, "LastLogin": account.LastLogin, "CrtDt": account.CrtDt}
	return accountCollection.InsertOne(context.TODO(), doc)
}

func (as *AccountService) UpdateLastLogin(mongoDB *mongo.Client, email string) (*mongo.UpdateResult, error) {
	accountCollection := mongoDB.Database("Account").Collection("AccountCollection")
	filter := bson.M{"_id": email}
	update := bson.M{"$set": bson.M{"LastLogin": time.Now()}}
	return accountCollection.UpdateOne(context.TODO(), filter, update)
}
