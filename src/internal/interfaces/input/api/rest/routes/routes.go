package routes

import (
	"go-project/src/internal/core/coreinterfaces"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(userHandler coreinterfaces.UserAPIHandler) http.Handler {
	r := chi.NewRouter()

	// LOGIN USER
	r.Post("/login", userHandler.LoginHandler)
	// AUTH USER
	r.Get("/auth/me", userHandler.AuthMeHandler)
	// GET ALL PRODUCTS
	r.Get("/allproducts", userHandler.AllProductsHandler)
	// GET SINGLE PRODUCT
	r.Get("/singleproduct", userHandler.GetSingleProduct)
	// GET ALL CATEGORIES
	r.Get("/get-categories", userHandler.CategoryHandler)
	// GET PRODUCTS OF GIVEN CATEGORIES
	r.Post("/get-products", userHandler.ProductHandler)

	return r
}
