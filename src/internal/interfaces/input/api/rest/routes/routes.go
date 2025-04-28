package routes

import (
	"go-project/src/internal/core/coreinterfaces"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(userHandler coreinterfaces.UserAPIHandler) http.Handler {
	r := chi.NewRouter()

	r.Post("/login", userHandler.LoginHandler)
	r.Get("/auth/me", userHandler.AuthMeHandler)

	return r
}
