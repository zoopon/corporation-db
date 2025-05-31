package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"corporation-db/internal/domain"
	"corporation-db/internal/infrastructure"
)

// CorporationUsecase handles corporation business logic
type CorporationUsecase struct {
	corporationRepo domain.CorporationRepository
	gbizClient      *infrastructure.GBizClient
}

// NewCorporationUsecase creates a new CorporationUsecase
func NewCorporationUsecase(corporationRepo domain.CorporationRepository, gbizClient *infrastructure.GBizClient) *CorporationUsecase {
	return &CorporationUsecase{
		corporationRepo: corporationRepo,
		gbizClient:      gbizClient,
	}
}

// GetAll retrieves all corporations
func (u *CorporationUsecase) GetAll(ctx context.Context) ([]*domain.Corporation, error) {
	return u.corporationRepo.GetAll(ctx)
}

// GetByID retrieves a corporation by ID
func (u *CorporationUsecase) GetByID(ctx context.Context, id int64) (*domain.Corporation, error) {
	return u.corporationRepo.GetByID(ctx, id)
}

// GetCorporationByCorporateNumber retrieves a corporation by corporate number
func (u *CorporationUsecase) GetByCorporateNumber(ctx context.Context, corporateNumber string) (*domain.Corporation, error) {
	// Validate corporate number format
	if len(corporateNumber) != 13 {
		return nil, domain.ErrInvalidCorporateNumber
	}

	return u.corporationRepo.GetByCorporateNumber(ctx, corporateNumber)
}

// GetCorporations retrieves corporations with filtering and pagination
func (u *CorporationUsecase) GetCorporations(ctx context.Context, filter domain.CorporationFilter) ([]*domain.Corporation, int64, error) {
	return u.corporationRepo.GetWithFilter(ctx, filter)
}

// ImportFromGBizINFO downloads and imports corporation data from gBizINFO
func (u *CorporationUsecase) ImportFromGBizINFO(ctx context.Context) error {
	log.Println("Starting gBizINFO data import...")

	// Download CSV ZIP file
	log.Println("Downloading basic corporation information from gBizINFO...")
	zipPath, err := u.gbizClient.DownloadBasicInfoCSV(ctx)
	if err != nil {
		return fmt.Errorf("failed to download CSV from gBizINFO: %w", err)
	}
	defer u.gbizClient.Cleanup(zipPath)

	log.Printf("Downloaded ZIP file: %s\n", zipPath)

	// Extract and parse CSV
	log.Println("Extracting and parsing CSV data...")
	corporations, err := u.gbizClient.ExtractAndParseCSV(zipPath)
	if err != nil {
		return fmt.Errorf("failed to extract and parse CSV: %w", err)
	}

	log.Printf("Parsed %d corporations from CSV\n", len(corporations))

	if len(corporations) == 0 {
		log.Println("No corporations to import")
		return nil
	}

	// Bulk upsert to database
	log.Println("Starting bulk upsert to database...")
	startTime := time.Now()

	err = u.corporationRepo.BulkUpsert(ctx, corporations)
	if err != nil {
		return fmt.Errorf("failed to bulk upsert corporations: %w", err)
	}

	duration := time.Since(startTime)
	log.Printf("Successfully imported %d corporations in %v\n", len(corporations), duration)
	log.Printf("Average: %.2f corporations/second\n", float64(len(corporations))/duration.Seconds())

	return nil
}

// ImportFromGBizINFOWithProgress downloads and imports corporation data with progress tracking
func (u *CorporationUsecase) ImportFromGBizINFOWithProgress(ctx context.Context, progressCallback func(stage string, progress float64)) error {
	if progressCallback != nil {
		progressCallback("Starting import", 0)
	}

	log.Println("Starting gBizINFO data import...")

	// Download CSV ZIP file
	if progressCallback != nil {
		progressCallback("Downloading data", 10)
	}

	log.Println("Downloading basic corporation information from gBizINFO...")
	zipPath, err := u.gbizClient.DownloadBasicInfoCSV(ctx)
	if err != nil {
		return fmt.Errorf("failed to download CSV from gBizINFO: %w", err)
	}
	defer u.gbizClient.Cleanup(zipPath)

	if progressCallback != nil {
		progressCallback("Download completed", 30)
	}

	log.Printf("Downloaded ZIP file: %s\n", zipPath)

	// Extract and parse CSV
	if progressCallback != nil {
		progressCallback("Parsing CSV data", 40)
	}

	log.Println("Extracting and parsing CSV data...")
	corporations, err := u.gbizClient.ExtractAndParseCSV(zipPath)
	if err != nil {
		return fmt.Errorf("failed to extract and parse CSV: %w", err)
	}

	if progressCallback != nil {
		progressCallback("CSV parsing completed", 70)
	}

	log.Printf("Parsed %d corporations from CSV\n", len(corporations))

	if len(corporations) == 0 {
		if progressCallback != nil {
			progressCallback("No data to import", 100)
		}
		log.Println("No corporations to import")
		return nil
	}

	// Bulk upsert to database
	if progressCallback != nil {
		progressCallback("Importing to database", 80)
	}

	log.Println("Starting bulk upsert to database...")
	startTime := time.Now()

	err = u.corporationRepo.BulkUpsert(ctx, corporations)
	if err != nil {
		return fmt.Errorf("failed to bulk upsert corporations: %w", err)
	}

	duration := time.Since(startTime)

	if progressCallback != nil {
		progressCallback("Import completed", 100)
	}

	log.Printf("Successfully imported %d corporations in %v\n", len(corporations), duration)
	log.Printf("Average: %.2f corporations/second\n", float64(len(corporations))/duration.Seconds())

	return nil
}

// GetImportStats returns statistics about the import process
func (u *CorporationUsecase) GetImportStats(ctx context.Context) (*ImportStats, error) {
	corporations, err := u.corporationRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get corporations for stats: %w", err)
	}

	stats := &ImportStats{
		TotalCorporations: len(corporations),
		LastImportTime:    time.Now(), // This should be tracked separately in a real implementation
	}

	// Count by status
	statusCounts := make(map[string]int)
	for _, corp := range corporations {
		statusCounts[corp.Status]++
	}
	stats.StatusCounts = statusCounts

	return stats, nil
}

// GetCorporationByCorporateNumber retrieves a corporation by corporate number
func (u *CorporationUsecase) GetCorporationByCorporateNumber(ctx context.Context, corporateNumber string) (*domain.Corporation, error) {
	// Validate corporate number format
	if len(corporateNumber) != 13 {
		return nil, domain.ErrInvalidCorporateNumber
	}

	corp, err := u.corporationRepo.GetByCorporateNumber(ctx, corporateNumber)
	if err != nil {
		return nil, err
	}
	if corp == nil {
		return nil, domain.ErrCorporationNotFound
	}

	return corp, nil
}

// ImportStats represents import statistics
type ImportStats struct {
	TotalCorporations int            `json:"total_corporations"`
	StatusCounts      map[string]int `json:"status_counts"`
	LastImportTime    time.Time      `json:"last_import_time"`
}

// ImportFromTestCSV imports corporation data from a local test CSV file
func (u *CorporationUsecase) ImportFromTestCSV(ctx context.Context, csvPath string) error {
	log.Printf("Starting test CSV import from: %s", csvPath)

	// Load and parse test CSV
	corporations, err := u.gbizClient.LoadTestCSVFile(csvPath)
	if err != nil {
		return fmt.Errorf("failed to load test CSV: %w", err)
	}

	log.Printf("Parsed %d corporations from test CSV", len(corporations))

	if len(corporations) == 0 {
		log.Println("No corporations to import")
		return nil
	}

	// Bulk upsert to database
	log.Println("Starting bulk upsert to database...")
	startTime := time.Now()

	err = u.corporationRepo.BulkUpsert(ctx, corporations)
	if err != nil {
		return fmt.Errorf("failed to bulk upsert corporations: %w", err)
	}

	duration := time.Since(startTime)
	log.Printf("Successfully imported %d corporations in %v", len(corporations), duration)
	log.Printf("Average: %.2f corporations/second", float64(len(corporations))/duration.Seconds())

	return nil
}

// ImportFromZIPFile imports corporation data from a ZIP file with progress tracking
func (u *CorporationUsecase) ImportFromZIPFile(ctx context.Context, zipPath string, progressCallback func(stage string, progress float64)) error {
	if progressCallback != nil {
		progressCallback("Starting import", 0)
	}

	log.Printf("Starting import from ZIP file: %s", zipPath)

	// Extract and parse CSV
	if progressCallback != nil {
		progressCallback("Parsing CSV data", 10)
	}

	log.Println("Extracting and parsing CSV data...")
	corporations, err := u.gbizClient.ExtractAndParseCSV(zipPath)
	if err != nil {
		return fmt.Errorf("failed to extract and parse CSV: %w", err)
	}

	if progressCallback != nil {
		progressCallback("CSV parsing completed", 50)
	}

	log.Printf("Parsed %d corporations from CSV", len(corporations))

	if len(corporations) == 0 {
		if progressCallback != nil {
			progressCallback("No data to import", 100)
		}
		log.Println("No corporations to import")
		return nil
	}

	// Bulk upsert to database
	if progressCallback != nil {
		progressCallback("Importing to database", 60)
	}

	log.Println("Starting bulk upsert to database...")
	startTime := time.Now()

	err = u.corporationRepo.BulkUpsert(ctx, corporations)
	if err != nil {
		return fmt.Errorf("failed to bulk upsert corporations: %w", err)
	}

	duration := time.Since(startTime)

	if progressCallback != nil {
		progressCallback("Import completed", 100)
	}

	log.Printf("Successfully imported %d corporations in %v", len(corporations), duration)
	log.Printf("Average: %.2f corporations/second", float64(len(corporations))/duration.Seconds())

	return nil
}
