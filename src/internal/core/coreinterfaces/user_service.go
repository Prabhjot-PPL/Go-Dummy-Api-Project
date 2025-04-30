package coreinterfaces

import (
	"context"
	"go-project/src/internal/adaptors/external/dummyapi"
	"go-project/src/internal/core/dto"
)

type Service interface {
	LoginUser(ctx context.Context, request dummyapi.UserCredentials) (dto.LoginResponse, error)
	GetUserByToken(ctx context.Context, token string) (dto.AuthResponse, error)
	GetAllProducts(ctx context.Context) ([]dummyapi.Product, error)
	GetProductById(ctx context.Context, id string) (dummyapi.Product, error)
	GetCategories(ctx context.Context) ([]string, error)
	GetProducts(ctx context.Context, req []string) ([]dummyapi.Product, error)
}
