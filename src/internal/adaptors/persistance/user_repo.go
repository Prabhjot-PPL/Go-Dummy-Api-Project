package persistance

import (
	"fmt"
	"go-project/src/internal/core/dto"
)

type UserRepo struct {
	db *Database
}

func NewUserRepo(d *Database) UserRepo {
	return UserRepo{db: d}
}

func (u *UserRepo) StoreUser(user dto.LoginResponse) error {

	// _, err := u.db.db.Exec("INSERT INTO user (id, username, email, first_name, last_name, gender, access_token, refresh_token, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", user.Id, user.Username, user.Email, user.FirstName, user.LastName, user.Gender, user.AccessToken, user.RefreshToken)
	_, err := u.db.db.Exec(`
	INSERT INTO "user" 
	(username, email, first_name, last_name, gender, access_token, refresh_token, created_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())`,
		user.Id, user.Username, user.Email, user.FirstName, user.LastName, user.Gender, user.AccessToken, user.RefreshToken,
	)

	if err != nil {
		fmt.Println("Error executing query:", err)
		return err
	}
	return nil
}
