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
)

type UserHandler struct {
	userService coreinterfaces.Service
}

func NewUserHandler(userService coreinterfaces.Service) coreinterfaces.UserAPIHandler {
	return &UserHandler{userService: userService}
}

// ------------------------------USER AUTH----------------------------------

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

	// Set the response header to indicate JSON content
	w.Header().Set("Content-Type", "application/json")

	// Encode the loginResponse as JSON and write it to the response body
	err = json.NewEncoder(w).Encode(loginResponse)
	// fmt.Println(loginResponse)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(allProducts)
	if err != nil {
		http.Error(w, "Failed to encode products", http.StatusInternalServerError)
	} else {
		log.Println("All Products fetched successfully!!!")
	}
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
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

	pkg.WriteResponse(w, respCategories)

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// err = json.NewEncoder(w).Encode(respCategories)
	// if err != nil {
	// 	http.Error(w, "Failed to encode Categories", http.StatusInternalServerError)
	// } else {
	// 	log.Println("All Categories fetched successfully!!!")
	// }
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		http.Error(w, "Failed to encode products", http.StatusInternalServerError)
	} else {
		log.Println("Products of give categories are fetched successfully!!!")
	}
}
