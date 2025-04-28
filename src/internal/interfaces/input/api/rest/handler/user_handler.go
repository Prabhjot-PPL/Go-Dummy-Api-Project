package userhandler

import (
	"context"
	"encoding/json"
	"go-project/src/internal/adaptors/external/dummyapi"
	"go-project/src/internal/core/coreinterfaces"
	"log"
	"net/http"
	"time"
)

type UserHandler struct {
	userService coreinterfaces.Service
}

func NewUserHandler(userService coreinterfaces.Service) coreinterfaces.UserAPIHandler {
	return &UserHandler{userService: userService}
}

func (u *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// time.Sleep(5 * time.Second) // DELAY

	var requestData dummyapi.ApiRequestData

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// SENDING RESPONSE to usecase.go
	loginResponse, err := u.userService.LoginUser(ctx, requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Println("Error in login response ")
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
		log.Println("Error encoding login response:", err)
		return
	}
}

func (u *UserHandler) AuthMeHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// Get the cookie
	cookie, err := r.Cookie("access_token")
	if err != nil {
		http.Error(w, "Access token missing", http.StatusUnauthorized)
		return
	}
	token := cookie.Value

	// Pass token to usecase
	userData, err := u.userService.GetUserByToken(ctx, token)
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
