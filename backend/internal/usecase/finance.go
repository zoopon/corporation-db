package usecase

import (
	"context"
	"fmt"
	"log"

	"corporation-db/internal/domain"
	"corporation-db/internal/infrastructure"
)

// FinanceUsecase handles finance-related business logic
type FinanceUsecase struct {
	financeRepo domain.FinanceRepository
	gbizClient  *infrastructure.GBizClient
}

// NewFinanceUsecase creates a new finance use case
func NewFinanceUsecase(financeRepo domain.FinanceRepository) *FinanceUsecase {
	return &FinanceUsecase{
		financeRepo: financeRepo,
		gbizClient:  infrastructure.NewGBizClient(),
	}
}

// ImportFromZIPFile imports finance data from a ZIP file
func (u *FinanceUsecase) ImportFromZIPFile(ctx context.Context, zipPath string, progressCallback func(stage string, progress float64)) error {
	if progressCallback != nil {
		progressCallback("Starting finance import", 0)
	}

	log.Printf("Starting finance data import from ZIP file: %s", zipPath)

	// Use GBizClient to extract and process finance CSV
	totalProcessed, err := u.gbizClient.ExtractAndProcessFinanceCSVStream(zipPath, func(finances []*domain.Finance) error {
		// Batch insert finances
		return u.financeRepo.CreateBatch(finances)
	})

	if err != nil {
		return fmt.Errorf("failed to process finance ZIP file: %w", err)
	}

	if progressCallback != nil {
		progressCallback("Finance import completed", 100)
	}

	log.Printf("Successfully imported %d finance records", totalProcessed)
	return nil
}

// GetImportStats returns finance import statistics
func (u *FinanceUsecase) GetImportStats(ctx context.Context) (*FinanceImportStats, error) {
	count, err := u.financeRepo.Count()
	if err != nil {
		return nil, fmt.Errorf("failed to get finance count: %w", err)
	}

	return &FinanceImportStats{
		TotalRecords: int(count),
	}, nil
}

// FinanceImportStats represents finance import statistics
type FinanceImportStats struct {
	TotalRecords int `json:"total_records"`
}
