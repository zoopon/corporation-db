package infrastructure

import (
	"context"
	"database/sql"
	"fmt"

	"corporation-db/internal/domain"
	"corporation-db/internal/infrastructure/db"

	"github.com/google/uuid"
)

type baseRepository struct {
	queries *db.Queries
}

// NewBaseRepository creates a new BaseRepository
func NewBaseRepository(database *sql.DB) domain.BaseRepository {
	return &baseRepository{
		queries: db.New(database),
	}
}

func (r *baseRepository) Create(ctx context.Context, base *domain.Base) (*domain.Base, error) {
	// Generate UUIDv7 if not set
	if base.ID == uuid.Nil {
		base.ID = domain.NewUUIDv7()
	}

	params := db.CreateBaseParams{
		ID:              base.ID,
		CorporationID:   base.CorporationID,
		CorporateNumber: base.CorporateNumber,
		BaseName:        stringPtrToNullString(base.BaseName),
		CountryCode:     base.CountryCode,
		PostalCode:      stringPtrToNullString(base.PostalCode),
		Location:        stringPtrToNullString(base.Location),
		PhoneNumber:     stringPtrToNullString(base.PhoneNumber),
		FaxNumber:       stringPtrToNullString(base.FaxNumber),
		DataObtainedAt:  base.DataObtainedAt,
		DataSourceUrl:   base.DataSourceURL,
		IsHeadOffice:    base.IsHeadOffice,
	}

	result, err := r.queries.CreateBase(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create base: %w", err)
	}

	return r.convertToBase(result), nil
}

func (r *baseRepository) CreateBatch(ctx context.Context, bases []*domain.Base) error {
	for _, base := range bases {
		// Generate UUIDv7 if not set
		if base.ID == uuid.Nil {
			base.ID = domain.NewUUIDv7()
		}

		params := db.CreateBaseBatchParams{
			ID:              base.ID,
			CorporationID:   base.CorporationID,
			CorporateNumber: base.CorporateNumber,
			BaseName:        stringPtrToNullString(base.BaseName),
			CountryCode:     base.CountryCode,
			PostalCode:      stringPtrToNullString(base.PostalCode),
			Location:        stringPtrToNullString(base.Location),
			PhoneNumber:     stringPtrToNullString(base.PhoneNumber),
			FaxNumber:       stringPtrToNullString(base.FaxNumber),
			DataObtainedAt:  base.DataObtainedAt,
			DataSourceUrl:   base.DataSourceURL,
			IsHeadOffice:    base.IsHeadOffice,
		}

		err := r.queries.CreateBaseBatch(ctx, params)
		if err != nil {
			return fmt.Errorf("failed to create base batch: %w", err)
		}
	}

	return nil
}

func (r *baseRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Base, error) {
	result, err := r.queries.GetBaseByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get base by ID: %w", err)
	}

	return r.convertToBase(result), nil
}

func (r *baseRepository) GetByCorporationID(ctx context.Context, corporationID uuid.UUID) ([]*domain.Base, error) {
	results, err := r.queries.GetBasesByCorporationID(ctx, corporationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bases by corporation ID: %w", err)
	}

	bases := make([]*domain.Base, len(results))
	for i, result := range results {
		bases[i] = r.convertToBase(result)
	}

	return bases, nil
}

func (r *baseRepository) GetByCorporateNumber(ctx context.Context, corporateNumber string) ([]*domain.Base, error) {
	results, err := r.queries.GetBasesByCorporateNumber(ctx, corporateNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get bases by corporate number: %w", err)
	}

	bases := make([]*domain.Base, len(results))
	for i, result := range results {
		bases[i] = r.convertToBase(result)
	}

	return bases, nil
}

func (r *baseRepository) GetHeadOfficeByCorporateNumber(ctx context.Context, corporateNumber string) (*domain.Base, error) {
	result, err := r.queries.GetHeadOfficeByCorporateNumber(ctx, corporateNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get head office by corporate number: %w", err)
	}

	return r.convertToBase(result), nil
}

func (r *baseRepository) Update(ctx context.Context, base *domain.Base) (*domain.Base, error) {
	params := db.UpdateBaseParams{
		ID:             base.ID,
		BaseName:       stringPtrToNullString(base.BaseName),
		CountryCode:    base.CountryCode,
		PostalCode:     stringPtrToNullString(base.PostalCode),
		Location:       stringPtrToNullString(base.Location),
		PhoneNumber:    stringPtrToNullString(base.PhoneNumber),
		FaxNumber:      stringPtrToNullString(base.FaxNumber),
		DataObtainedAt: base.DataObtainedAt,
		DataSourceUrl:  base.DataSourceURL,
		IsHeadOffice:   base.IsHeadOffice,
	}

	result, err := r.queries.UpdateBase(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update base: %w", err)
	}

	return r.convertToBase(result), nil
}

func (r *baseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteBase(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete base: %w", err)
	}

	return nil
}

func (r *baseRepository) DeleteByCorporationID(ctx context.Context, corporationID uuid.UUID) error {
	err := r.queries.DeleteBasesByCorporationID(ctx, corporationID)
	if err != nil {
		return fmt.Errorf("failed to delete bases by corporation ID: %w", err)
	}

	return nil
}

func (r *baseRepository) ListAll(ctx context.Context, limit, offset int) ([]*domain.Base, error) {
	params := db.ListAllBasesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	results, err := r.queries.ListAllBases(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list all bases: %w", err)
	}

	bases := make([]*domain.Base, len(results))
	for i, result := range results {
		bases[i] = r.convertToBase(result)
	}

	return bases, nil
}

func (r *baseRepository) SearchByName(ctx context.Context, name string, limit, offset int) ([]*domain.Base, error) {
	params := db.SearchBasesByNameParams{
		Column1: sql.NullString{String: name, Valid: name != ""},
		Limit:   int32(limit),
		Offset:  int32(offset),
	}

	results, err := r.queries.SearchBasesByName(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to search bases by name: %w", err)
	}

	bases := make([]*domain.Base, len(results))
	for i, result := range results {
		bases[i] = r.convertToBase(result)
	}

	return bases, nil
}

func (r *baseRepository) GetByCountry(ctx context.Context, countryCode string, limit, offset int) ([]*domain.Base, error) {
	params := db.GetBasesByCountryParams{
		CountryCode: countryCode,
		Limit:       int32(limit),
		Offset:      int32(offset),
	}

	results, err := r.queries.GetBasesByCountry(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get bases by country: %w", err)
	}

	bases := make([]*domain.Base, len(results))
	for i, result := range results {
		bases[i] = r.convertToBase(result)
	}

	return bases, nil
}

func (r *baseRepository) CountByCorporateNumber(ctx context.Context, corporateNumber string) (int64, error) {
	count, err := r.queries.CountBasesByCorporateNumber(ctx, corporateNumber)
	if err != nil {
		return 0, fmt.Errorf("failed to count bases by corporate number: %w", err)
	}

	return count, nil
}

// convertToBase converts db.Basis to domain.Base
func (r *baseRepository) convertToBase(dbBase db.Basis) *domain.Base {
	return &domain.Base{
		ID:              dbBase.ID,
		CorporationID:   dbBase.CorporationID,
		CorporateNumber: dbBase.CorporateNumber,
		BaseName:        nullStringToPtr(dbBase.BaseName),
		CountryCode:     dbBase.CountryCode,
		PostalCode:      nullStringToPtr(dbBase.PostalCode),
		Location:        nullStringToPtr(dbBase.Location),
		PhoneNumber:     nullStringToPtr(dbBase.PhoneNumber),
		FaxNumber:       nullStringToPtr(dbBase.FaxNumber),
		DataObtainedAt:  dbBase.DataObtainedAt,
		DataSourceURL:   dbBase.DataSourceUrl,
		IsHeadOffice:    dbBase.IsHeadOffice,
		CreatedAt:       dbBase.CreatedAt,
		UpdatedAt:       dbBase.UpdatedAt,
	}
}

// nullStringToPtr converts sql.NullString to *string
func nullStringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

// stringPtrToNullString converts *string to sql.NullString
func stringPtrToNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{Valid: false}
}
