package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"corporation-db/internal/infrastructure"
	"corporation-db/internal/presentation"
	"corporation-db/internal/usecase"
	"corporation-db/internal/utils"

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

	// Initialize repositories
	corporationRepo := infrastructure.NewCorporationRepository(db)
	baseRepo := infrastructure.NewBaseRepository(db)

	// Initialize gBiz client (for batch operations)
	gbizClient := infrastructure.NewGBizClient()

	// Initialize text converter
	textConverter := utils.NewTextConverter()

	// Initialize use cases
	corporationUsecase := usecase.NewCorporationUsecase(corporationRepo, gbizClient, textConverter)
	baseUsecase := usecase.NewBaseUsecase(baseRepo, corporationRepo)

	// Set base repository in corporation usecase to avoid circular dependency
	corporationUsecase.SetBaseRepo(baseRepo)

	// Initialize router with handlers
	router := presentation.NewRouter(corporationUsecase, baseUsecase)
	r := router.SetupRoutes()

	// Server
	port := getEnv("PORT", "8080")
	addr := fmt.Sprintf(":%s", port)

	log.Printf("🚀 Development server starting on port %s with hot reload", port)
	log.Printf("✅ Health check available at http://localhost:%s/health", port)
	log.Printf("🏢 Corporations API available at http://localhost:%s/corporations", port)
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
