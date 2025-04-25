package user

import (
	"go-project/src/internal/core/dto"
)

type User struct {
	Uid      int
	Username string `json:"username"`
	Password string `json:"password"`
}

type Service interface {
	LoginUser(request User) (dto.LoginResponse, error)
	GetUserByToken(token string) (dto.AuthResponse, error)
}
