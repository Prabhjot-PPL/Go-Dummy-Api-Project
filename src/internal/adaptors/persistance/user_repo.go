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
	// Check if user with the same email already exists
	var count int
	err := u.db.db.QueryRow(`
		SELECT COUNT(*) FROM "user" WHERE email = $1 OR username = $2
	`, user.Email, user.Username).Scan(&count)
	if err != nil {
		fmt.Println("Error checking for existing user:", err)
		return err
	}

	// If a user already exists with the same email or username, return an error
	if count > 0 {
		return fmt.Errorf("user with this email or username already exists")
	}

	// Proceed with inserting the new user
	_, err = u.db.db.Exec(`
		INSERT INTO "user" 
		(username, email, first_name, last_name, gender, access_token, refresh_token, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())`,
		user.Username, user.Email, user.FirstName, user.LastName, user.Gender, user.AccessToken, user.RefreshToken,
	)

	if err != nil {
		fmt.Println("Error executing query:", err)
		return err
	}
	return nil
}
