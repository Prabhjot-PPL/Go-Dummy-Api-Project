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
		protected.Get("/products", userHandler.AllProductsHandler)
		// GET SINGLE PRODUCT
		protected.Get("/products/{id}", userHandler.GetSingleProduct)
		// GET ALL CATEGORIES
		protected.Get("/categories", userHandler.CategoryHandler)
		// GET PRODUCTS OF GIVEN CATEGORIES
		protected.Post("/categories/{cat_id}/products", userHandler.ProductHandler)
		// UPDATE PRODUCT
		protected.Put("/products/{id}", userHandler.UpdateProductHandler)
		// DELETE PRODUCT
		protected.Delete("/products/{id}", userHandler.DeleteProductHandler)

	})

	return r
}

/*

3rd party taking too long to respond (e.g. 10 mins)



// apis

user first api - response with "your request has been accepted"
second call - "still processing unitl we fetched data from 3rd party"
.
.
nth call - "responsd with success / failure (3rd party)"
*/
