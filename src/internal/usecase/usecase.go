package usecase

import (
	"encoding/json"
	"fmt"
	"go-project/src/internal/adaptors/external"
	"go-project/src/internal/adaptors/persistance"
	"go-project/src/internal/core/dto"
	"go-project/src/internal/core/user"
)

type UserService struct {
	userRepo persistance.UserRepo
}

func NewUserService(userRepo persistance.UserRepo) UserService {
	return UserService{userRepo: userRepo}
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

	fmt.Println("idhar se shuru")
	fmt.Printf("%+v", loginResp)

	return loginResp, nil
}
