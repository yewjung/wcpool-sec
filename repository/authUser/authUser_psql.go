package authUserRepo

import (
	"database/sql"
	"sec/models"
)

type AuthUserRepo struct{}

func (authRepo *AuthUserRepo) UserExist(db *sql.DB, email string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM authuser WHERE email=$1", email).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

// create new user
func (authRepo *AuthUserRepo) CreateUser(db *sql.DB, user *models.AuthUser) error {
	_, err := db.Exec("INSERT INTO authuser(email, passwordhash) VALUES($1, $2)", user.Email, user.PasswordHash)
	if err != nil {
		return err
	}
	return nil
}

// update user
func (authRepo *AuthUserRepo) UpdateUser(db *sql.DB, user *models.AuthUser) error {
	_, err := db.Exec("UPDATE authuser SET passwordhash=$1, userid=$2 WHERE email=$3", user.PasswordHash, user.UserId, user.Email)
	if err != nil {
		return err
	}
	return nil
}
