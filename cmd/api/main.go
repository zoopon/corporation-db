package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"corporation-db/internal/infrastructure"
	"corporation-db/internal/presentation"
	"corporation-db/internal/usecase"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Database connection
	dbConfig := infrastructure.NewDBConfig()
	db, err := dbConfig.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Repositories
	userRepo := infrastructure.NewUserRepository(db)

	// Use cases
	userUsecase := usecase.NewUserUsecase(userRepo)

	// Router
	router := presentation.NewRouter(userUsecase)
	r := router.SetupRoutes()

	// Server
	port := getEnv("PORT", "8080")
	addr := fmt.Sprintf(":%s", port)

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
