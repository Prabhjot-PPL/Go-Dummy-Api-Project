package dummyapi

// use context here..
// remove db from here (UserRepo struct)
// get data from struct
// change struct name from userRepo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ApiInterface interface {
	GetUser(ctx context.Context, requestData UserCredentials) (*http.Response, error)
	GetUserByToken(ctx context.Context, token string) (*http.Response, error)
	GetAllProducts(ctx context.Context) ([]Product, error)
	GetProductById(ctx context.Context, id string) (Product, error)
	GetProductCategories(ctx context.Context) ([]string, error)
	GetProductsByCategory(ctx context.Context, category string) ([]Product, error)
	UpdateProduct(ctx context.Context, id string, updateData Product) (Product, error)
	DeleteProduct(ctx context.Context, id string) (Product, error)
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

// ------------------------------USER AUTH----------------------------------

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

// ------------------------------ PRODUCT ----------------------------------

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

// GET ALL PRODUCTS
func (key ApiImplementation) GetAllProducts(ctx context.Context) ([]Product, error) {

	url := key.baseUrl + "/products"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var products ProductsResponse
	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		log.Print("Error decoding products : ", err)
		return nil, err
	}

	// fmt.Println(products)

	return products.Products, err
}

// GET PRODUCT BY ID
func (key ApiImplementation) GetProductById(ctx context.Context, id string) (Product, error) {

	url := key.baseUrl + "/products/" + id

	// fmt.Println(url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Println("Error creating request for GetProductById")
		return Product{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error making request for GetProductById")
		return Product{}, err
	}

	defer resp.Body.Close()

	var product Product
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		return Product{}, err
	}

	return product, nil
}

// GET CATEGORIES
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

// GET PRODUCTS BASED ON CATEGORIES GIVEN
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

// UPDATE PRODUCT
func (key ApiImplementation) UpdateProduct(ctx context.Context, id string, updateData Product) (Product, error) {
	url := key.baseUrl + "/products/" + id

	jsonBody, err := json.Marshal(updateData)
	if err != nil {
		return Product{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return Product{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Product{}, err
	}
	defer resp.Body.Close()

	fmt.Println(resp.Body)

	var updatedProduct Product
	if err := json.NewDecoder(resp.Body).Decode(&updatedProduct); err != nil {
		return Product{}, err
	}

	fmt.Println(updatedProduct)

	return updatedProduct, nil
}

// DELETE PRODUCT
func (key ApiImplementation) DeleteProduct(ctx context.Context, id string) (Product, error) {
	url := key.baseUrl + "/products/" + id

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		log.Println("Error creating request for DeleteProduct")
		return Product{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error executing request for DeleteProduct")
		return Product{}, err
	}
	defer resp.Body.Close()

	var product Product
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		return Product{}, err
	}

	return product, nil
}
