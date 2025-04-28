package dummyapi

// use context here..
// remove db from here (UserRepo struct)
// get data from struct
// change struct name from userRepo

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type ApiImplementation struct {
	BaseUrl string
}

type ApiRequestData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (key ApiImplementation) GetUser(ctx context.Context, requestData ApiRequestData) (*http.Response, error) {

	url := key.BaseUrl + "/auth/login"

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

	url := key.BaseUrl + "/auth/me"

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
