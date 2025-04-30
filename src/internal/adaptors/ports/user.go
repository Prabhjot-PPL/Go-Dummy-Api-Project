package ports

import (
	"context"
	"go-project/src/internal/core/dto"
)

type UserRepository interface {
	CheckUserExist(ctx context.Context, username string) error
	StoreUser(ctx context.Context, user dto.LoginResponse) error
}
