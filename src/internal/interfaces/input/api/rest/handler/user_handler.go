package userhandler

import (
	"encoding/json"
	"fmt"
	"go-project/src/internal/core/user"
	"net/http"
)

type UserHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) UserHandler {
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

	// Set the access_token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    loginResponse.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	// Set the response header to indicate JSON content
	w.Header().Set("Content-Type", "application/json")

	// Encode the loginResponse as JSON and write it to the response body
	err = json.NewEncoder(w).Encode(loginResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error encoding response"))
		fmt.Println("Error encoding login response:", err)
		return
	}

	fmt.Println()
	// fmt.Println("login Response : ", loginResponse)
}

func (u *UserHandler) AuthMeHandler(w http.ResponseWriter, r *http.Request) {

	// Get the cookie
	cookie, err := r.Cookie("access_token")
	if err != nil {
		http.Error(w, "Access token missing", http.StatusUnauthorized)
		return
	}
	token := cookie.Value

	// Pass token to usecase
	userData, err := u.userService.GetUserByToken(token)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Respond with user data
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userData)
	if err != nil {
		http.Error(w, "Failed to encode user data", http.StatusInternalServerError)
		return
	}
}
