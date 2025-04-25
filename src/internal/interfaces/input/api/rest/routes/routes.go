package routes

import (
	userhandler "go-project/src/internal/interfaces/input/api/rest/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(userHandler *userhandler.UserHandler) http.Handler {
	r := chi.NewRouter()

	r.Post("/login", userHandler.LoginHandler)

	return r
}
