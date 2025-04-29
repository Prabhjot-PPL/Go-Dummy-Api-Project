package dummyapi

// use context here..
// remove db from here (UserRepo struct)
// get data from struct
// change struct name from userRepo

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type ApiInterface interface {
}

type ApiImplementation struct {
	baseUrl string
}

func New(baseUrl string) ApiInterface {
	return &ApiImplementation{
		baseUrl: baseUrl,
	}
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (key ApiImplementation) GetUser(ctx context.Context, requestData UserCredentials) (*http.Response, error) {

	url := key.baseUrl + "/auth/login"

	// Encode the request body
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	// Create a new HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (key ApiImplementation) GetUserByToken(ctx context.Context, token string) (*http.Response, error) {

	url := key.baseUrl + "/auth/me"

	// Create a new GET request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set Authorization header with Bearer token
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the request using http.DefaultClient
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (key ApiImplementation) GetProductCategories(ctx context.Context) ([]string, error) {

	url := key.baseUrl + "/products/category-list"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var categories []string
	err = json.NewDecoder(resp.Body).Decode(&categories)
	if err != nil {
		return nil, err
	}

	return categories, nil

}

type Product struct {
	Id          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	Discount    float64 `json:"discountPercentage"`
	Rating      float64 `json:"rating"`
}

type ProductsResponse struct {
	Products []Product `json:"products"`
}

func (key ApiImplementation) GetProductsByCategory(ctx context.Context, category string) ([]Product, error) {

	url := key.baseUrl + "/products/category/" + category

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var productsResp ProductsResponse
	err = json.NewDecoder(resp.Body).Decode(&productsResp)
	if err != nil {
		log.Println("Error decoding response:", err)
		return nil, err
	}

	return productsResp.Products, nil
}
