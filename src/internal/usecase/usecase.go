package usecase

import (
	"context"
	"encoding/json"
	"go-project/src/internal/adaptors/external/dummyapi"
	"go-project/src/internal/config"
	"go-project/src/internal/core/coreinterfaces"
	"go-project/src/internal/core/dto"
	"log"
)

type userService struct {
	userRepo coreinterfaces.UserRepository
	keys     dummyapi.ApiImplementation
}

func NewUserService(userRepo coreinterfaces.UserRepository) coreinterfaces.Service {
	return &userService{userRepo: userRepo}
}

// LOGIN request to dummy_api
func (u *userService) LoginUser(ctx context.Context, requestData dummyapi.ApiRequestData) (dto.LoginResponse, error) {

	config := config.LoadConfig()
	baseUrl := dummyapi.ApiImplementation{BaseUrl: config.Dummy_API}

	// CHECK if user already exist in DB or not
	err := u.userRepo.CheckUserExist(ctx, requestData.Username)
	if err != nil {
		log.Println("Error checking user existence:", err)
		return dto.LoginResponse{}, err
	}

	// SEND username, password to dummy api
	resp, err := baseUrl.GetUser(ctx, requestData)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	// closing response body
	defer resp.Body.Close()

	// fetching RESPONSE from dummy_api
	var loginResp dto.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	// STORING response data to DB
	err = u.userRepo.StoreUser(ctx, loginResp)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return loginResp, nil
}

// AUTH user
func (u *userService) GetUserByToken(ctx context.Context, token string) (dto.AuthResponse, error) {

	config := config.LoadConfig()
	baseUrl := dummyapi.ApiImplementation{BaseUrl: config.Dummy_API}

	resp, err := baseUrl.GetUserByToken(ctx, token)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	defer resp.Body.Close()

	var userResp dto.AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&userResp)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return userResp, nil
}
