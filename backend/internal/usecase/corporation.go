package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"corporation-db/internal/domain"
	"corporation-db/internal/infrastructure"
	"corporation-db/internal/utils"

	"github.com/google/uuid"
)

// CorporationUsecase handles corporation business logic
type CorporationUsecase struct {
	corporationRepo domain.CorporationRepository
	baseRepo        domain.BaseRepository
	gbizClient      *infrastructure.GBizClient
	textConverter   *utils.TextConverter
}

// NewCorporationUsecase creates a new CorporationUsecase
func NewCorporationUsecase(corporationRepo domain.CorporationRepository, gbizClient *infrastructure.GBizClient, textConverter *utils.TextConverter) *CorporationUsecase {
	return &CorporationUsecase{
		corporationRepo: corporationRepo,
		baseRepo:        nil, // Will be set later via SetBaseRepo
		gbizClient:      gbizClient,
		textConverter:   textConverter,
	}
}

// SetBaseRepo sets the base repository (used to avoid circular dependency)
func (u *CorporationUsecase) SetBaseRepo(baseRepo domain.BaseRepository) {
	u.baseRepo = baseRepo
}

// GetAll retrieves all corporations
func (u *CorporationUsecase) GetAll(ctx context.Context) ([]*domain.Corporation, error) {
	return u.corporationRepo.GetAll(ctx)
}

// GetByID retrieves a corporation by ID
func (u *CorporationUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Corporation, error) {
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
	// Normalize name for search if provided
	if filter.Name != nil && *filter.Name != "" {
		normalizedName := u.textConverter.NormalizeForSearch(*filter.Name)
		filter.Name = &normalizedName
	}

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

	// Use streaming import for memory efficiency
	return u.ImportFromZIPFileStream(ctx, zipPath, nil)
}

// ImportFromGBizINFOStream downloads and imports corporation data using streaming processing for memory efficiency
func (u *CorporationUsecase) ImportFromGBizINFOStream(ctx context.Context, progressCallback func(stage string, progress float64)) error {
	if progressCallback != nil {
		progressCallback("Starting import", 0)
	}

	log.Println("Starting gBizINFO data import with streaming...")

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

	// Use streaming import for memory efficiency
	return u.ImportFromZIPFileStream(ctx, zipPath, progressCallback)
}

// ImportFromGBizINFOWithProgress downloads and imports corporation data with progress tracking
func (u *CorporationUsecase) ImportFromGBizINFOWithProgress(ctx context.Context, progressCallback func(stage string, progress float64)) error {
	// Use the new streaming method for memory efficiency
	return u.ImportFromGBizINFOStream(ctx, progressCallback)
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
	// Use the new streaming method for memory efficiency
	return u.ImportFromZIPFileStream(ctx, zipPath, progressCallback)
}

// ImportFromZIPFileStream imports corporation data from a ZIP file using streaming processing for memory efficiency
func (u *CorporationUsecase) ImportFromZIPFileStream(ctx context.Context, zipPath string, progressCallback func(stage string, progress float64)) error {
	if progressCallback != nil {
		progressCallback("Starting import", 0)
	}

	log.Printf("Starting streaming import from ZIP file: %s", zipPath)

	// Extract and process CSV in streaming batches
	if progressCallback != nil {
		progressCallback("Processing CSV data", 10)
	}

	log.Println("Extracting and processing CSV data in streaming batches...")
	startTime := time.Now()

	totalProcessed := 0
	batchCount := 0

	// Define batch processor function
	batchProcessor := func(batch []*domain.CreateCorporationRequest) error {
		batchCount++
		batchStartTime := time.Now()

		// Process this batch to database
		err := u.corporationRepo.BulkUpsert(ctx, batch)
		if err != nil {
			return fmt.Errorf("failed to upsert batch %d: %w", batchCount, err)
		}

		// Create head office records for new corporations if base repository is available
		if u.baseRepo != nil {
			u.createHeadOfficesForBatch(ctx, batch)
		}

		batchDuration := time.Since(batchStartTime)
		totalProcessed += len(batch)

		log.Printf("Processed batch %d: %d records in %v (total: %d)",
			batchCount, len(batch), batchDuration, totalProcessed)

		// Update progress based on processed records (rough estimate)
		if progressCallback != nil {
			// Progress from 10% to 90% during processing
			// Use a more conservative estimate for large files
			progress := 10.0 + (float64(totalProcessed)/5000000.0)*80.0 // Assume up to 5M records for progress calculation
			if progress > 90.0 {
				progress = 90.0
			}
			progressCallback("Processing batches", progress)
		}

		return nil
	}

	// Use streaming CSV processing
	processed, err := u.gbizClient.ExtractAndProcessCSVStream(zipPath, batchProcessor)
	if err != nil {
		return fmt.Errorf("failed to extract and process CSV stream: %w", err)
	}

	duration := time.Since(startTime)

	if progressCallback != nil {
		progressCallback("Import completed", 100)
	}

	if processed == 0 {
		log.Println("No corporations to import")
		return nil
	}

	log.Printf("Successfully imported %d corporations in %d batches over %v",
		processed, batchCount, duration)
	log.Printf("Average: %.2f corporations/second", float64(processed)/duration.Seconds())

	return nil
}

// createHeadOfficesForBatch creates head office records for a batch of corporations
func (u *CorporationUsecase) createHeadOfficesForBatch(ctx context.Context, batch []*domain.CreateCorporationRequest) {
	var headOffices []*domain.Base

	for _, corpReq := range batch {
		// Get the corporation to get its ID
		corp, err := u.corporationRepo.GetByCorporateNumber(ctx, corpReq.CorporateNumber)
		if err != nil {
			log.Printf("Warning: Could not find corporation %s for head office creation: %v", corpReq.CorporateNumber, err)
			continue
		}

		// Check if head office already exists
		existing, err := u.baseRepo.GetHeadOfficeByCorporateNumber(ctx, corpReq.CorporateNumber)
		if err == nil && existing != nil {
			continue // Head office already exists
		}

		// Create head office base from corporation data
		headOffice := domain.NewHeadOfficeBase(corp.ID, corp.CorporateNumber)
		headOffice.CountryCode = "JP" // Default to Japan

		// Build address from corporation data
		address := ""
		if corp.PostalCode != nil && *corp.PostalCode != "" {
			headOffice.PostalCode = corp.PostalCode
			address = *corp.PostalCode + " "
		}
		if corp.Location != nil && *corp.Location != "" {
			address += *corp.Location
		}
		if address != "" {
			headOffice.Location = &address
		} else {
			unknown := "住所不明" // Unknown address
			headOffice.Location = &unknown
		}

		headOffice.DataObtainedAt = corp.CreatedAt
		headOffice.DataSourceURL = "https://info.gbiz.go.jp/"

		headOffices = append(headOffices, headOffice)
	}

	// Batch create head offices
	if len(headOffices) > 0 {
		err := u.baseRepo.CreateBatch(ctx, headOffices)
		if err != nil {
			log.Printf("Warning: Failed to create head offices batch: %v", err)
		} else {
			log.Printf("Created %d head office records", len(headOffices))
		}
	}
}
