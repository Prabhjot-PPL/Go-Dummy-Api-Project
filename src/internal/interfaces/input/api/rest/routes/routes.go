package routes

import (
	"go-project/src/internal/adaptors/ports"
	"go-project/src/internal/core/coreinterfaces"
	"go-project/src/internal/interfaces/input/api/rest/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(userHandler coreinterfaces.UserAPIHandler, userRepo ports.UserRepository) http.Handler {
	r := chi.NewRouter()
	authMiddleware := middleware.NewAuthMiddleware(userRepo)

	// LOGIN USER
	r.Post("/login", userHandler.LoginHandler)
	// AUTH USER
	r.Get("/auth/me", userHandler.AuthMeHandler)

	r.Group(func(protected chi.Router) {
		protected.Use(authMiddleware.Middleware)

		// GET ALL PRODUCTS
		protected.Get("/allproducts", userHandler.AllProductsHandler)
		// GET SINGLE PRODUCT
		protected.Get("/singleproduct", userHandler.GetSingleProduct)
		// GET ALL CATEGORIES
		protected.Get("/get-categories", userHandler.CategoryHandler)
		// GET PRODUCTS OF GIVEN CATEGORIES
		protected.Post("/get-products", userHandler.ProductHandler)
		// UPDATE PRODUCT
		protected.Put("/updateproduct/{id}", userHandler.UpdateProductHandler)
		// DELETE PRODUCT
		protected.Delete("/deleteproduct/{id}", userHandler.DeleteProductHandler)

	})

	return r
}
