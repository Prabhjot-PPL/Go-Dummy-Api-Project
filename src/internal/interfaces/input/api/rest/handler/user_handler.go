package userhandler

import (
	"context"
	"encoding/json"
	"go-project/src/internal/adaptors/external/dummyapi"
	"go-project/src/internal/core/coreinterfaces"
	"go-project/src/pkg"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService coreinterfaces.Service
}

func NewUserHandler(userService coreinterfaces.Service) coreinterfaces.UserAPIHandler {
	return &UserHandler{userService: userService}
}

// ------------------------------USER AUTH----------------------------------

// LOGIN USER
func (u *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()

	// time.Sleep(5 * time.Second) // DELAY

	var requestData dummyapi.UserCredentials

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

	errString := "Failed to login user."
	successSting := "User logged-in successfully!!!"
	pkg.WriteResponse(w, loginResponse, errString, successSting)
}

// AUTH USER
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

	errString := "Failed to authenticate user."
	successSting := "User authenticated successfully!!!"
	pkg.WriteResponse(w, userData, errString, successSting)
}

// ------------------------------ PRODUCT ----------------------------------

// GET ALL PRODUCTS
func (u *UserHandler) AllProductsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	allProducts, err := u.userService.GetAllProducts(ctx)
	if err != nil {
		http.Error(w, "Failed to get products", http.StatusInternalServerError)
		return
	}

	errString := "Failed to fetch products user."
	successSting := "Products fetched successfully!!!"
	pkg.WriteResponse(w, allProducts, errString, successSting)
}

// GET SINGLE PRODUCT
func (u *UserHandler) GetSingleProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing product ID", http.StatusBadRequest)
		return
	}

	product, err := u.userService.GetProductById(ctx, id)
	if err != nil {
		http.Error(w, "Failed to fetch product", http.StatusInternalServerError)
		return
	}

	errString := "Failed to get product by id."
	successSting := "Successfully fetched Product by Id!!!"
	pkg.WriteResponse(w, product, errString, successSting)
}

// GET ALL CATEGORIES
func (u *UserHandler) CategoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()

	respCategories, err := u.userService.GetCategories(ctx)
	if err != nil {
		http.Error(w, "Failed to get categories", http.StatusInternalServerError)
		return
	}

	errString := "Failed to encode Categories"
	successSting := "All Categories fetched successfully!!!"
	pkg.WriteResponse(w, respCategories, errString, successSting)
}

// GET PRODUCTS OF GIVEN CATEGORIES
func (u *UserHandler) ProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()

	var reqBody struct {
		Categories []string `json:"categories"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	products, err := u.userService.GetProducts(ctx, reqBody.Categories)
	if err != nil {
		http.Error(w, "Failed to get products", http.StatusInternalServerError)
		return
	}

	errString := "Failed to fetch products of given categories."
	successSting := "Successfully fetched Products of given categories!!!"
	pkg.WriteResponse(w, products, errString, successSting)
}

// UPDATE PRODUCT
func (u *UserHandler) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing product ID", http.StatusBadRequest)
		return
	}

	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedProduct, err := u.userService.UpdateProduct(ctx, id, updateData)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedProduct)
}

// DELETE PRODUCT
func (u *UserHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing product ID", http.StatusBadRequest)
		return
	}

	product, err := u.userService.DeleteProduct(ctx, id)
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}
