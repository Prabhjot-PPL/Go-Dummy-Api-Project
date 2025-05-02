package usecase

import (
	"context"
	"encoding/json"
	"go-project/src/internal/adaptors/external/dummyapi"
	"go-project/src/internal/adaptors/ports"
	"go-project/src/internal/core/coreinterfaces"
	"go-project/src/internal/core/dto"
	"log"
	"sync"
)

type userService struct {
	userRepo ports.UserRepository
	keys     dummyapi.ApiInterface
}

func NewUserService(userRepo ports.UserRepository, baseurl dummyapi.ApiInterface) coreinterfaces.Service {
	return &userService{
		userRepo: userRepo,
		keys:     baseurl,
	}
}

// ------------------------------USER AUTH----------------------------------

// LOGIN request to dummy_api
func (u *userService) LoginUser(ctx context.Context, requestData dummyapi.UserCredentials) (dto.LoginResponse, error) {

	// CHECK if user already exist in DB or not
	err := u.userRepo.CheckUserExist(ctx, requestData.Username)
	if err != nil {
		log.Println("Error checking user existence:", err)
		return dto.LoginResponse{}, err
	}

	// SEND username, password to dummy api
	resp, err := u.keys.GetUser(ctx, requestData)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	// closing response body
	defer resp.Body.Close()

	// fetching RESPONSE from dummy_api
	var loginResp dto.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	// STORING response data to DB
	err = u.userRepo.StoreUser(ctx, loginResp)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	// fmt.Println(loginResp)

	return loginResp, nil
}

// AUTH user
func (u *userService) GetUserByToken(ctx context.Context, token string) (dto.AuthResponse, error) {

	resp, err := u.keys.GetUserByToken(ctx, token)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	defer resp.Body.Close()

	var userResp dto.AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&userResp)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return userResp, nil
}

// ------------------------------ PRODUCT ----------------------------------

// GET ALL PRODUCTS
func (u *userService) GetAllProducts(ctx context.Context) ([]dummyapi.Product, error) {

	resp, err := u.keys.GetAllProducts(ctx)
	if err != nil {
		log.Println("Error getting products : ", err)
		return resp, err
	}

	return resp, err
}

// GET PRODUCT BY ID
func (u *userService) GetProductById(ctx context.Context, id string) (dummyapi.Product, error) {
	resp, err := u.keys.GetProductById(ctx, id)
	if err != nil {
		log.Println("Error getting product by id : ", err)
		return dummyapi.Product{}, err
	}

	return resp, err
}

// GET ALL CATEGORIES
func (u *userService) GetCategories(ctx context.Context) ([]string, error) {

	categories, err := u.keys.GetProductCategories(ctx)
	if err != nil {
		log.Println("Error getting product category ", err)
		return nil, err
	}

	// log.Println("Fetched categories : ", categories)

	return categories, err
}

// GET PRODUCTS OF GIVEN CATEGORIES
func (u *userService) GetProducts(ctx context.Context, categories []string) ([]dummyapi.Product, error) {

	var wg sync.WaitGroup
	var mu sync.Mutex
	var allProducts []dummyapi.Product

	for _, category := range categories {
		wg.Add(1)
		go func(cat string) {
			defer wg.Done()
			products, err := u.keys.GetProductsByCategory(ctx, cat)
			if err != nil {
				return
			}
			mu.Lock()
			allProducts = append(allProducts, products...)
			mu.Unlock()
		}(category)
	}

	wg.Wait()

	return allProducts, nil
}

// UPDATE PRODUCT
func (u *userService) UpdateProduct(ctx context.Context, id string, updateData dummyapi.Product) (dummyapi.Product, error) {
	resp, err := u.keys.UpdateProduct(ctx, id, updateData)
	if err != nil {
		log.Println("Error deleting product by id: ", err)
		return dummyapi.Product{}, err
	}
	return resp, nil
}

// DELETE PRODUCT
func (u *userService) DeleteProduct(ctx context.Context, id string) (dummyapi.Product, error) {
	resp, err := u.keys.DeleteProduct(ctx, id)
	if err != nil {
		log.Println("Error deleting product by id: ", err)
		return dummyapi.Product{}, err
	}
	return resp, nil
}
