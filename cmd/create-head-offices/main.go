package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"corporation-db/internal/domain"
	"corporation-db/internal/infrastructure"
	"corporation-db/internal/usecase"

	_ "github.com/lib/pq"
)

func main() {
	// Get database URL from environment variable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/corporation_db?sslmode=disable"
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	corporationRepo := infrastructure.NewCorporationRepository(db)
	baseRepo := infrastructure.NewBaseRepository(db)

	// Initialize usecases
	baseUsecase := usecase.NewBaseUsecase(baseRepo, corporationRepo)

	ctx := context.Background()

	// Get corporations without head office
	query := `
		SELECT c.id, c.corporate_number, c.name, c.postal_code, c.location, c.created_at 
		FROM corporations c 
		WHERE c.corporate_number NOT IN (
			SELECT DISTINCT corporate_number FROM bases WHERE is_head_office = true
		) 
		LIMIT 10
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Fatalf("Failed to query corporations: %v", err)
	}
	defer rows.Close()

	created := 0
	for rows.Next() {
		corp := &domain.Corporation{}
		err := rows.Scan(
			&corp.ID,
			&corp.CorporateNumber,
			&corp.Name,
			&corp.PostalCode,
			&corp.Location,
			&corp.CreatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan corporation: %v", err)
			continue
		}

		// Create head office for this corporation
		_, err = baseUsecase.CreateHeadOfficeFromCorporation(ctx, corp)
		if err != nil {
			log.Printf("Failed to create head office for %s (%s): %v", corp.Name, corp.CorporateNumber, err)
			continue
		}

		created++
		log.Printf("Created head office for %s (%s)", corp.Name, corp.CorporateNumber)
	}

	log.Printf("Successfully created %d head offices", created)
}
