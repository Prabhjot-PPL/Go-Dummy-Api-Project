package main

import (
	"go-project/src/internal/adaptors/persistance"
	userhandler "go-project/src/internal/interfaces/input/api/rest/handler"
	"go-project/src/internal/interfaces/input/api/rest/routes"
	"go-project/src/internal/usecase"
	"log"
	"net/http"
)

func main() {

	database, err := persistance.ConnectToDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database : %v", err)
	}

	UserRepo := persistance.NewUserRepo(database)
	UserService := usecase.NewUserService(UserRepo)
	userHandler := userhandler.NewUserHandler(UserService)

	router := routes.InitRoutes(userHandler)

	err = http.ListenAndServe(":8080", router)

	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	log.Println("Server running on http://localhost:8080")
}
