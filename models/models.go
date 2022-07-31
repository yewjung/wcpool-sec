package models

import (
	"database/sql"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthUser struct {
	Email        string `json:"email"`
	PasswordHash string `json:"passwordhash"`
}

type Account struct {
	Email      string          `bson:"_id"`
	Username   string          `bson:"Username"`
	Parties    map[string]bool `bson:"Parties"`
	Permission map[string]bool `bson:"Permission"`
	LastLogin  time.Time       `bson:"LastLogin"`
	CrtDt      time.Time       `bson:"CrtDt"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type Error struct {
	ErrorMessage string `json:"errorMessage"`
}

type Storage struct {
	PostgresUserDB    *sql.DB
	MongoAccountDB    *mongo.Client
	RedisAccountCache *redis.Client
}

type VerificationMethod func(account Account, partyid string) bool
