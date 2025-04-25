package usecase

import (
	"encoding/json"
	"go-project/src/internal/adaptors/external"
	"go-project/src/internal/adaptors/persistance"
	"go-project/src/internal/core/dto"
	"go-project/src/internal/core/user"
)

type UserService struct {
	userRepo persistance.UserRepo
}

func NewUserService(userRepo persistance.UserRepo) user.Service {
	return &UserService{userRepo: userRepo}
}

// type LoginResponse struct {
// 	Id           int
// 	Username     string
// 	Email        string
// 	FirstName    string
// 	LastName     string
// 	Gender       string
// 	AccessToken  string
// 	RefreshToken string
// }

func (u *UserService) LoginUser(requestData user.User) (dto.LoginResponse, error) {

	resp, err := external.GetUser(requestData.Username, requestData.Password)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	defer resp.Body.Close()

	var loginResp dto.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	err = u.userRepo.StoreUser(loginResp)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return loginResp, nil
}

func (u *UserService) GetUserByToken(token string) (dto.AuthResponse, error) {

	resp, err := external.GetUserByToken(token)
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
