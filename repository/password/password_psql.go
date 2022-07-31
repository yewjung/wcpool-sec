package passwordRepo

import (
	"database/sql"
	"sec/models"
)

type PasswordRepo struct {
	DB *sql.DB
}

func (passwordRepo *PasswordRepo) UserExist(email string) bool {
	var count int
	err := passwordRepo.DB.QueryRow("SELECT COUNT(*) FROM password WHERE email=$1", email).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

// create new user
func (passwordRepo *PasswordRepo) CreateUser(user *models.AuthUser) error {
	_, err := passwordRepo.DB.Exec("INSERT INTO password(email, passwordhash) VALUES($1, $2)", user.Email, user.PasswordHash)
	if err != nil {
		return err
	}
	return nil
}

// update user
func (passwordRepo *PasswordRepo) UpdateUser(user *models.AuthUser) error {
	_, err := passwordRepo.DB.Exec("UPDATE password SET passwordhash=$1 WHERE email=$2", user.PasswordHash, user.Email)
	if err != nil {
		return err
	}
	return nil
}

func (passwordRepo *PasswordRepo) GetUser(email string) (models.AuthUser, error) {
	row := passwordRepo.DB.QueryRow("SELECT email, passwordhash FROM password WHERE email=$1", email)
	authUser := models.AuthUser{}
	err := row.Scan(&authUser.Email, &authUser.PasswordHash)
	return authUser, err
}
