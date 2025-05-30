package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"corporation-db/internal/domain"
	"corporation-db/internal/infrastructure/db"

	_ "github.com/lib/pq"
)

// CorporationRepository implements domain.CorporationRepository
type CorporationRepository struct {
	db      *sql.DB
	queries *db.Queries
}

// NewCorporationRepository creates a new CorporationRepository
func NewCorporationRepository(database *sql.DB) domain.CorporationRepository {
	return &CorporationRepository{
		db:      database,
		queries: db.New(database),
	}
}

// GetAll retrieves all corporations
func (r *CorporationRepository) GetAll(ctx context.Context) ([]*domain.Corporation, error) {
	dbCorporations, err := r.queries.GetCorporations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get corporations: %w", err)
	}

	corporations := make([]*domain.Corporation, len(dbCorporations))
	for i, dbCorp := range dbCorporations {
		corp, err := r.dbToDomainCorporation(&dbCorp)
		if err != nil {
			return nil, fmt.Errorf("failed to convert corporation %d: %w", dbCorp.ID, err)
		}
		corporations[i] = corp
	}

	return corporations, nil
}

// GetWithFilter retrieves corporations with filtering and pagination
func (r *CorporationRepository) GetWithFilter(ctx context.Context, filter domain.CorporationFilter) ([]*domain.Corporation, int64, error) {
	// Prepare filter parameters - use empty string for null values
	corporateNumber := ""
	name := ""
	prefecture := ""
	status := ""
	corporateType := ""

	if filter.CorporateNumber != nil {
		corporateNumber = *filter.CorporateNumber
	}
	if filter.Name != nil {
		name = *filter.Name
	}
	if filter.Prefecture != nil {
		prefecture = *filter.Prefecture
	}
	if filter.Status != nil {
		status = *filter.Status
	}
	if filter.CorporateType != nil {
		corporateType = *filter.CorporateType
	}

	// Get total count
	total, err := r.queries.CountCorporationsWithFilter(ctx, db.CountCorporationsWithFilterParams{
		Column1: corporateNumber,
		Column2: name,
		Column3: prefecture,
		Column4: status,
		Column5: corporateType,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count corporations: %w", err)
	}

	// Get filtered corporations
	dbCorporations, err := r.queries.GetCorporationsWithFilter(ctx, db.GetCorporationsWithFilterParams{
		Column1: corporateNumber,
		Column2: name,
		Column3: prefecture,
		Column4: status,
		Column5: corporateType,
		Limit:   int32(filter.Limit),
		Offset:  int32(filter.Offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get filtered corporations: %w", err)
	}

	corporations := make([]*domain.Corporation, len(dbCorporations))
	for i, dbCorp := range dbCorporations {
		corp, err := r.dbToDomainCorporation(&dbCorp)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to convert corporation %d: %w", dbCorp.ID, err)
		}
		corporations[i] = corp
	}

	return corporations, total, nil
}

// GetByID retrieves a corporation by ID
func (r *CorporationRepository) GetByID(ctx context.Context, id int64) (*domain.Corporation, error) {
	dbCorp, err := r.queries.GetCorporationByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrCorporationNotFound
		}
		return nil, fmt.Errorf("failed to get corporation by ID: %w", err)
	}

	return r.dbToDomainCorporation(&dbCorp)
}

// GetByCorporateNumber retrieves a corporation by corporate number
func (r *CorporationRepository) GetByCorporateNumber(ctx context.Context, corporateNumber string) (*domain.Corporation, error) {
	dbCorp, err := r.queries.GetCorporationByCorporateNumber(ctx, corporateNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrCorporationNotFound
		}
		return nil, fmt.Errorf("failed to get corporation by corporate number: %w", err)
	}

	return r.dbToDomainCorporation(&dbCorp)
}

// BulkUpsert performs bulk upsert of corporations
func (r *CorporationRepository) BulkUpsert(ctx context.Context, corporations []*domain.CreateCorporationRequest) error {
	if len(corporations) == 0 {
		return nil
	}

	// Use transaction for bulk operations
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	txQueries := r.queries.WithTx(tx)

	// Process corporations in batches to avoid memory issues
	batchSize := 1000
	for i := 0; i < len(corporations); i += batchSize {
		end := i + batchSize
		if end > len(corporations) {
			end = len(corporations)
		}

		batch := corporations[i:end]

		for j, corp := range batch {
			params := db.UpsertCorporationParams{
				CorporateNumber:     corp.CorporateNumber,
				Name:                corp.Name,
				NameKana:            r.stringToNullString(corp.NameKana),
				EnglishName:         r.stringToNullString(corp.EnglishName),
				PostalCode:          r.stringToNullString(corp.PostalCode),
				Address:             r.stringToNullString(corp.Address),
				PrefectureCode:      r.stringToNullString(corp.PrefectureCode),
				CityCode:            r.stringToNullString(corp.CityCode),
				FoundingDate:        r.timeToNullTime(corp.FoundingDate),
				Status:              corp.Status,
				CorporateType:       r.stringToNullString(corp.CorporateType),
				CapitalStock:        r.int64ToNullInt64(corp.CapitalStock),
				EmployeeNumber:      r.int32ToNullInt32(corp.EmployeeNumber),
				Representative:      r.stringToNullString(corp.Representative),
				BusinessDescription: r.stringToNullString(corp.BusinessDescription),
				Industry:            r.stringToNullString(corp.Industry),
				Website:             r.stringToNullString(corp.Website),
				Phone:               r.stringToNullString(corp.Phone),
				Email:               r.stringToNullString(corp.Email),
				LastUpdated:         r.timeToNullTime(corp.LastUpdated),
			}

			_, err := txQueries.UpsertCorporation(ctx, params)
			if err != nil {
				return fmt.Errorf("failed to upsert corporation %s (batch %d, item %d): %w",
					corp.CorporateNumber, i/batchSize+1, j+1, err)
			}
		}

		fmt.Printf("Processed batch %d/%d (%d corporations)\n",
			i/batchSize+1, (len(corporations)+batchSize-1)/batchSize, len(batch))
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Printf("Successfully upserted %d corporations\n", len(corporations))
	return nil
}

// Helper function to convert database Corporation to domain Corporation
func (r *CorporationRepository) dbToDomainCorporation(dbCorp *db.Corporation) (*domain.Corporation, error) {
	corp := &domain.Corporation{
		ID:                  int64(dbCorp.ID),
		CorporateNumber:     dbCorp.CorporateNumber,
		Name:                dbCorp.Name,
		NameKana:            r.nullStringToString(dbCorp.NameKana),
		EnglishName:         r.nullStringToString(dbCorp.EnglishName),
		PostalCode:          r.nullStringToString(dbCorp.PostalCode),
		Address:             r.nullStringToString(dbCorp.Address),
		PrefectureCode:      r.nullStringToString(dbCorp.PrefectureCode),
		CityCode:            r.nullStringToString(dbCorp.CityCode),
		FoundingDate:        r.nullTimeToTime(dbCorp.FoundingDate),
		Status:              dbCorp.Status,
		CorporateType:       r.nullStringToString(dbCorp.CorporateType),
		CapitalStock:        r.nullInt64ToInt64(dbCorp.CapitalStock),
		EmployeeNumber:      r.nullInt32ToInt32(dbCorp.EmployeeNumber),
		Representative:      r.nullStringToString(dbCorp.Representative),
		BusinessDescription: r.nullStringToString(dbCorp.BusinessDescription),
		Industry:            r.nullStringToString(dbCorp.Industry),
		Website:             r.nullStringToString(dbCorp.Website),
		Phone:               r.nullStringToString(dbCorp.Phone),
		Email:               r.nullStringToString(dbCorp.Email),
		LastUpdated:         r.nullTimeToTime(dbCorp.LastUpdated),
		CreatedAt:           r.nullTimeToTimeValue(dbCorp.CreatedAt),
		UpdatedAt:           r.nullTimeToTimeValue(dbCorp.UpdatedAt),
	}

	return corp, nil
}

// Helper functions for null type conversions
func (r *CorporationRepository) stringToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func (r *CorporationRepository) nullStringToString(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

func (r *CorporationRepository) timeToNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func (r *CorporationRepository) nullTimeToTime(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

func (r *CorporationRepository) int64ToNullInt64(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: *i, Valid: true}
}

func (r *CorporationRepository) nullInt64ToInt64(ni sql.NullInt64) *int64 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int64
}

func (r *CorporationRepository) int32ToNullInt32(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: *i, Valid: true}
}

func (r *CorporationRepository) nullInt32ToInt32(ni sql.NullInt32) *int32 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int32
}

func (r *CorporationRepository) nullTimeToTimeValue(nt sql.NullTime) time.Time {
	if !nt.Valid {
		return time.Time{}
	}
	return nt.Time
}
