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
	r.Get("/get-categories", userHandler.CategoryHandler)
	r.Post("/get-products", userHandler.ProductHandler)

	return r
}
