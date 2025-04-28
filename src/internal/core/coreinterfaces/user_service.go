package coreinterfaces

import (
	"context"
	"go-project/src/internal/adaptors/external/dummyapi"
	"go-project/src/internal/core/dto"
)

type Service interface {
	LoginUser(ctx context.Context, request dummyapi.ApiRequestData) (dto.LoginResponse, error)
	GetUserByToken(ctx context.Context, token string) (dto.AuthResponse, error)
}
