package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sec/models"
	"time"

	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountService struct {
	MongoDB *mongo.Client
	Cache   *redis.Client
}

func (as *AccountService) CreateAccount(account models.Account) (*mongo.InsertOneResult, error) {
	accountCollection := as.MongoDB.Database("Account").Collection("AccountCollection")
	doc := bson.M{
		"_id":        account.Email,
		"Username":   account.Username,
		"LastLogin":  account.LastLogin,
		"CrtDt":      account.CrtDt,
		"Parties":    make(map[string]bool, 0),
		"Permission": make(map[string]bool, 0),
	}
	return accountCollection.InsertOne(context.TODO(), doc)
}

func (as *AccountService) UpdateLastLogin(email string) (*mongo.UpdateResult, error) {
	accountCollection := as.MongoDB.Database("Account").Collection("AccountCollection")
	filter := bson.M{"_id": email}
	update := bson.M{"$set": bson.M{"LastLogin": time.Now()}}
	return accountCollection.UpdateOne(context.TODO(), filter, update)
}

func (as *AccountService) IsUserFromParty(account models.Account, partyid string) bool {
	_, ok := account.Parties[partyid]
	return ok
}

func (as *AccountService) IsUserAdminOfParty(account models.Account, partyid string) bool {
	permissionKey := constructPermissionKey(partyid, "admin")
	_, ok := account.Permission[permissionKey]
	return ok
}

func (as *AccountService) FindByEmail(email string) models.Account {
	cacheResult, err := as.Cache.Get(context.Background(), email).Result()
	if err == nil {
		cachedAccount := models.Account{}
		json.Unmarshal([]byte(cacheResult), &cachedAccount)
		return cachedAccount
	}
	accountCollection := as.MongoDB.Database("Account").Collection("AccountCollection")
	filter := bson.M{"_id": email}
	result := accountCollection.FindOne(context.Background(), filter)
	account := models.Account{}
	err = result.Decode(&account)
	if err != nil {
		log.Panic(err)
		return account
	}
	accountByte, err := json.Marshal(account)
	if err != nil {
		return account
	}
	as.Cache.Set(context.Background(), email, accountByte, 0)
	return account
}

func constructPermissionKey(partyid string, userGroup string) string {
	return fmt.Sprintf("%s$%s", partyid, userGroup)
}
