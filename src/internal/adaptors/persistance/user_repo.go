package persistance

import (
	"context"
	"fmt"
	"go-project/src/internal/adaptors/ports"
	"go-project/src/internal/core/dto"
	"log"
)

type UserRepo struct {
	db *Database
}

func NewUserRepo(d *Database) ports.UserRepository {
	return &UserRepo{db: d}
}

func (u *UserRepo) CheckUserExist(ctx context.Context, username string) error {

	// Check if user with the same email already exists
	var count int
	err := u.db.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM "userinfo" WHERE username = $1
	`, username).Scan(&count)
	if err != nil {
		log.Println("Error checking for existing user:", err)
		return err
	}

	// If a user already exists with the same email or username, return an error
	if count > 0 {
		return fmt.Errorf("user with this username already exists")
	}

	return nil
}

// INSERT user data in DB
func (u *UserRepo) StoreUser(ctx context.Context, user dto.LoginResponse) error {
	_, err := u.db.db.ExecContext(ctx, `
		INSERT INTO "userinfo" 
		(username, email, first_name, last_name, gender, access_token, refresh_token, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
	`,
		user.Username, user.Email, user.FirstName, user.LastName, user.Gender, user.AccessToken, user.RefreshToken,
	)

	if err != nil {
		log.Println("Error executing insert query:", err)
		return err
	}
	return nil
}

// AUTH USER (for middleware)
func (u *UserRepo) IsTokenValid(ctx context.Context, token string) (bool, error) {

	// Validate access_token
	var count int
	err := u.db.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM "userinfo" WHERE access_token = $1
	`, token).Scan(&count)
	if err != nil {
		log.Println("Error validating access_token : ", err)
		return false, err
	}

	// If a access_token doesn't exist
	if count == 0 {
		return false, fmt.Errorf("user with this username already exists")
	}

	return true, nil
}
