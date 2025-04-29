package coreinterfaces

import (
	"context"
	"go-project/src/internal/adaptors/external/dummyapi"
	"go-project/src/internal/core/dto"
)

type Service interface {
	LoginUser(ctx context.Context, request dummyapi.UserCredentials) (dto.LoginResponse, error)
	GetUserByToken(ctx context.Context, token string) (dto.AuthResponse, error)
	GetCategories(ctx context.Context) error
	GetProducts(ctx context.Context, req []string) ([]dummyapi.Product, error)
}
