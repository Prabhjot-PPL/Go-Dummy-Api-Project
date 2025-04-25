package userhandler

import (
	"encoding/json"
	"fmt"
	"go-project/src/internal/core/user"
	"go-project/src/internal/usecase"
	"net/http"
)

type UserHandler struct {
	userService usecase.UserService
}

func NewUserHandler(userService usecase.UserService) UserHandler {
	return UserHandler{userService: userService}
}

func (u *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var requestData user.User

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// fmt.Println("Request Data : ", requestData)

	loginResponse, err := u.userService.LoginUser(requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		fmt.Println("Error in login response ")
		return
	}

	// w.Write()

	fmt.Println()
	fmt.Println("login Response : ", loginResponse)
}
