package usecase

import (
	"context"
	"fmt"

	"corporation-db/internal/domain"

	"github.com/google/uuid"
)

// BaseUsecase handles base/branch office business logic
type BaseUsecase struct {
	baseRepo        domain.BaseRepository
	corporationRepo domain.CorporationRepository
}

// NewBaseUsecase creates a new BaseUsecase
func NewBaseUsecase(baseRepo domain.BaseRepository, corporationRepo domain.CorporationRepository) *BaseUsecase {
	return &BaseUsecase{
		baseRepo:        baseRepo,
		corporationRepo: corporationRepo,
	}
}

// CreateBase creates a new base/branch office
func (u *BaseUsecase) CreateBase(ctx context.Context, base *domain.Base) (*domain.Base, error) {
	// Validate that the corporation exists
	_, err := u.corporationRepo.GetByID(ctx, base.CorporationID)
	if err != nil {
		return nil, fmt.Errorf("corporation not found: %w", err)
	}

	// Validate that corporate number matches
	corp, err := u.corporationRepo.GetByCorporateNumber(ctx, base.CorporateNumber)
	if err != nil {
		return nil, fmt.Errorf("corporation with corporate number %s not found: %w", base.CorporateNumber, err)
	}
	if corp.ID != base.CorporationID {
		return nil, fmt.Errorf("corporation ID mismatch: expected %s, got %s", corp.ID, base.CorporationID)
	}

	// If this is being set as head office, check if one already exists
	if base.IsHeadOffice {
		existingHeadOffice, err := u.baseRepo.GetHeadOfficeByCorporateNumber(ctx, base.CorporateNumber)
		if err == nil && existingHeadOffice != nil {
			return nil, fmt.Errorf("head office already exists for corporation %s", base.CorporateNumber)
		}
	}

	return u.baseRepo.Create(ctx, base)
}

// GetBaseByID retrieves a base by ID
func (u *BaseUsecase) GetBaseByID(ctx context.Context, id uuid.UUID) (*domain.Base, error) {
	return u.baseRepo.GetByID(ctx, id)
}

// GetBasesByCorporationID retrieves all bases for a corporation
func (u *BaseUsecase) GetBasesByCorporationID(ctx context.Context, corporationID uuid.UUID) ([]*domain.Base, error) {
	return u.baseRepo.GetByCorporationID(ctx, corporationID)
}

// GetBasesByCorporateNumber retrieves all bases for a corporate number
func (u *BaseUsecase) GetBasesByCorporateNumber(ctx context.Context, corporateNumber string) ([]*domain.Base, error) {
	// Validate corporate number format
	if len(corporateNumber) != 13 {
		return nil, domain.ErrInvalidCorporateNumber
	}

	return u.baseRepo.GetByCorporateNumber(ctx, corporateNumber)
}

// GetHeadOfficeByCorporateNumber retrieves the head office for a corporate number
func (u *BaseUsecase) GetHeadOfficeByCorporateNumber(ctx context.Context, corporateNumber string) (*domain.Base, error) {
	// Validate corporate number format
	if len(corporateNumber) != 13 {
		return nil, domain.ErrInvalidCorporateNumber
	}

	return u.baseRepo.GetHeadOfficeByCorporateNumber(ctx, corporateNumber)
}

// UpdateBase updates a base
func (u *BaseUsecase) UpdateBase(ctx context.Context, base *domain.Base) (*domain.Base, error) {
	// Validate that the base exists
	existing, err := u.baseRepo.GetByID(ctx, base.ID)
	if err != nil {
		return nil, fmt.Errorf("base not found: %w", err)
	}

	// If changing to head office, check if one already exists for this corporation
	if base.IsHeadOffice && !existing.IsHeadOffice {
		existingHeadOffice, err := u.baseRepo.GetHeadOfficeByCorporateNumber(ctx, base.CorporateNumber)
		if err == nil && existingHeadOffice != nil && existingHeadOffice.ID != base.ID {
			return nil, fmt.Errorf("head office already exists for corporation %s", base.CorporateNumber)
		}
	}

	return u.baseRepo.Update(ctx, base)
}

// DeleteBase deletes a base by ID
func (u *BaseUsecase) DeleteBase(ctx context.Context, id uuid.UUID) error {
	// Validate that the base exists
	_, err := u.baseRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("base not found: %w", err)
	}

	return u.baseRepo.Delete(ctx, id)
}

// DeleteBasesByCorporationID deletes all bases for a corporation
func (u *BaseUsecase) DeleteBasesByCorporationID(ctx context.Context, corporationID uuid.UUID) error {
	return u.baseRepo.DeleteByCorporationID(ctx, corporationID)
}

// ListBases retrieves bases with pagination
func (u *BaseUsecase) ListBases(ctx context.Context, limit, offset int) ([]*domain.Base, error) {
	if limit <= 0 {
		limit = 50 // Default limit
	}
	if limit > 1000 {
		limit = 1000 // Maximum limit
	}
	if offset < 0 {
		offset = 0
	}

	return u.baseRepo.ListAll(ctx, limit, offset)
}

// SearchBasesByName searches bases by name with pagination
func (u *BaseUsecase) SearchBasesByName(ctx context.Context, name string, limit, offset int) ([]*domain.Base, error) {
	if name == "" {
		return nil, fmt.Errorf("search name cannot be empty")
	}

	if limit <= 0 {
		limit = 50 // Default limit
	}
	if limit > 1000 {
		limit = 1000 // Maximum limit
	}
	if offset < 0 {
		offset = 0
	}

	return u.baseRepo.SearchByName(ctx, name, limit, offset)
}

// GetBasesByCountry retrieves bases by country code with pagination
func (u *BaseUsecase) GetBasesByCountry(ctx context.Context, countryCode string, limit, offset int) ([]*domain.Base, error) {
	if countryCode == "" {
		return nil, fmt.Errorf("country code cannot be empty")
	}

	if limit <= 0 {
		limit = 50 // Default limit
	}
	if limit > 1000 {
		limit = 1000 // Maximum limit
	}
	if offset < 0 {
		offset = 0
	}

	return u.baseRepo.GetByCountry(ctx, countryCode, limit, offset)
}

// CreateHeadOfficeFromCorporation creates a head office base from corporation data
func (u *BaseUsecase) CreateHeadOfficeFromCorporation(ctx context.Context, corp *domain.Corporation) (*domain.Base, error) {
	// Check if head office already exists
	existingHeadOffice, err := u.baseRepo.GetHeadOfficeByCorporateNumber(ctx, corp.CorporateNumber)
	if err == nil && existingHeadOffice != nil {
		return existingHeadOffice, nil // Head office already exists
	}

	// Create new head office base from corporation data
	headOffice := domain.NewHeadOfficeBase(corp.ID, corp.CorporateNumber)
	headOffice.CountryCode = "JP" // Default to Japan for now

	// Build address from postal code and location
	address := ""
	if corp.PostalCode != nil && *corp.PostalCode != "" {
		address = *corp.PostalCode + " "
		headOffice.PostalCode = corp.PostalCode
	}
	if corp.Location != nil {
		address += *corp.Location
	}
	headOffice.Location = &address

	headOffice.DataObtainedAt = corp.CreatedAt
	headOffice.DataSourceURL = "https://gbiz-info.go.jp/"

	return u.baseRepo.Create(ctx, headOffice)
}
