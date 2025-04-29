package usecase

import (
	"context"
	"encoding/json"
	"go-project/src/internal/adaptors/external/dummyapi"
	"go-project/src/internal/config"
	"go-project/src/internal/core/coreinterfaces"
	"go-project/src/internal/core/dto"
	"log"
	"sync"
)

type userService struct {
	userRepo coreinterfaces.UserRepository
	keys     dummyapi.ApiInterface
}

func NewUserService(userRepo coreinterfaces.UserRepository) coreinterfaces.Service {
	return &userService{userRepo: userRepo}
}

// LOGIN request to dummy_api
func (u *userService) LoginUser(ctx context.Context, requestData dummyapi.UserCredentials) (dto.LoginResponse, error) {

	config := config.LoadConfig()
	// baseUrl := dummyapi.ApiImplementation{BaseUrl: config.Dummy_API}
	u.keys.BaseUrl = config.Dummy_API

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

	config := config.LoadConfig()
	// baseUrl := dummyapi.ApiImplementation{BaseUrl: config.Dummy_API}
	u.keys.BaseUrl = "gfh"

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

// --------------------------PRODUCT-----------------------------

func (u *userService) GetCategories(ctx context.Context) error {

	config := config.LoadConfig()
	u.keys.BaseUrl = config.Dummy_API

	categories, err := u.keys.GetProductCategories(ctx)
	if err != nil {
		log.Println("Error getting product category ", err)
		return err
	}

	log.Println("Fetched categories : ", categories)

	return err
}

func (u *userService) GetProducts(ctx context.Context, categories []string) ([]dummyapi.Product, error) {
	config := config.LoadConfig()
	u.keys.BaseUrl = config.Dummy_API

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
