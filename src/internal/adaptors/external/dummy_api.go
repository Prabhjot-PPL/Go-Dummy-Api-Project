package external

// use context here..
// remove db from here (UserRepo struct)
// get data from struct
// change struct name from userRepo

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// type UserRepo struct {
// 	db *persistance.Database
// }

// func NewUserRepo(d *persistance.Database) UserRepo {
// 	return UserRepo{db: d}
// }

func GetUser(username string, password string) (*http.Response, error) {

	url := "https://dummyjson.com/auth/login"

	// Prepare the JSON body
	requestBody := map[string]interface{}{
		"username":      "emilys",
		"password":      "emilyspass",
		"expiresInMins": 30,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		panic(err)
	}

	// Send POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	return resp, nil
}
