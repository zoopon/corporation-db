package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"corporation-db/internal/domain"
	"corporation-db/internal/infrastructure"

	_ "github.com/lib/pq"
)

func main() {
	// Connect to database
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "password")
	dbName := getEnvOrDefault("DB_NAME", "corporation_db")
	dbSSLMode := getEnvOrDefault("DB_SSL_MODE", "disable")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Successfully connected to database")

	// Create repository
	financeRepo := infrastructure.NewFinanceRepository(db)

	// Test 1: Create finance using NewFinance constructor (should have UUIDv7)
	fmt.Println("\n=== Test 1: Creating Finance with NewFinance constructor ===")

	finance1 := domain.NewFinance("1111111111111") // Use existing corporate number
	finance1.CorporateName = "Test Finance Record 1"
	finance1.BusinessYear = "2024"
	finance1.SalesRevenue = "1000000"

	fmt.Printf("Finance ID (before save): %s\n", finance1.ID.String())
	fmt.Printf("Finance ID version: %d\n", finance1.ID.Version())

	if err := financeRepo.Create(finance1); err != nil {
		log.Printf("Failed to create finance1: %v", err)
	} else {
		fmt.Printf("Successfully created finance1 with ID: %s\n", finance1.ID.String())
	}

	// Test 2: Create finance manually then save (should get UUIDv7 from repository)
	fmt.Println("\n=== Test 2: Creating Finance manually (repository should add UUIDv7) ===")

	finance2 := &domain.Finance{
		CorporateNumber: "2222222222222", // Use existing corporate number
		CorporateName:   "Test Finance Record 2",
		BusinessYear:    "2024",
		SalesRevenue:    "2000000",
	}

	fmt.Printf("Finance ID (before save): %s\n", finance2.ID.String())

	if err := financeRepo.Create(finance2); err != nil {
		log.Printf("Failed to create finance2: %v", err)
	} else {
		fmt.Printf("Successfully created finance2 with ID: %s\n", finance2.ID.String())
	}

	// Test 3: Create batch of finances using NewFinance (all should have UUIDv7)
	fmt.Println("\n=== Test 3: Creating batch of Finance records ===")

	financesBatch := []*domain.Finance{
		domain.NewFinance("3333333333333"),
		domain.NewFinance("1111111111111"),
		domain.NewFinance("2222222222222"),
	}

	for i, f := range financesBatch {
		f.CorporateName = fmt.Sprintf("Batch Finance Record %d", i+1)
		f.BusinessYear = "2024"
		f.SalesRevenue = fmt.Sprintf("%d000000", i+3)
		fmt.Printf("Batch Finance %d ID: %s (version: %d)\n", i+1, f.ID.String(), f.ID.Version())
	}

	if err := financeRepo.CreateBatch(financesBatch); err != nil {
		log.Printf("Failed to create finance batch: %v", err)
	} else {
		fmt.Printf("Successfully created batch of %d finance records\n", len(financesBatch))
	}

	// Test 4: Verify all records are in database with proper UUIDs
	fmt.Println("\n=== Test 4: Verifying records in database ===")

	count, err := financeRepo.Count()
	if err != nil {
		log.Printf("Failed to count finances: %v", err)
	} else {
		fmt.Printf("Total finance records in database: %d\n", count)
	}

	// Check records for each corporate number
	for _, corpNum := range []string{"1111111111111", "2222222222222", "3333333333333"} {
		finances, err := financeRepo.GetByCorporateNumber(corpNum)
		if err != nil {
			log.Printf("Failed to get finances for %s: %v", corpNum, err)
			continue
		}

		fmt.Printf("\nFinance records for corporate number %s:\n", corpNum)
		for i, f := range finances {
			fmt.Printf("  Record %d: ID=%s (version=%d), Name=%s, Revenue=%s\n",
				i+1, f.ID.String(), f.ID.Version(), f.CorporateName, f.SalesRevenue)
		}
	}

	fmt.Println("\n=== UUIDv7 Generation Test Completed ===")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
