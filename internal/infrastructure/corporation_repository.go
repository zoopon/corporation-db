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
		corporations[i] = r.convertToDomain(&dbCorp)
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
	prefectureCode := ""

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
	if filter.PrefectureCode != nil {
		prefectureCode = *filter.PrefectureCode
	}

	// Get total count
	total, err := r.queries.CountCorporationsWithFilter(ctx, db.CountCorporationsWithFilterParams{
		Column1: corporateNumber,
		Column2: name,
		Column3: prefecture,
		Column4: status,
		Column5: prefectureCode,
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
		Column5: prefectureCode,
		Limit:   int32(filter.Limit),
		Offset:  int32(filter.Offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get filtered corporations: %w", err)
	}

	corporations := make([]*domain.Corporation, len(dbCorporations))
	for i, dbCorp := range dbCorporations {
		corporations[i] = r.convertToDomain(&dbCorp)
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

	return r.convertToDomain(&dbCorp), nil
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

	return r.convertToDomain(&dbCorp), nil
}

// BulkUpsert performs bulk upsert of corporations with transaction commits every 1000 records
func (r *CorporationRepository) BulkUpsert(ctx context.Context, corporations []*domain.CreateCorporationRequest) error {
	if len(corporations) == 0 {
		return nil
	}

	// Process corporations in batches with individual transactions for each batch
	// Reduced batch size for better memory efficiency on 4GB systems
	batchSize := 500
	totalBatches := (len(corporations) + batchSize - 1) / batchSize

	for i := 0; i < len(corporations); i += batchSize {
		end := i + batchSize
		if end > len(corporations) {
			end = len(corporations)
		}

		batch := corporations[i:end]
		batchNum := i/batchSize + 1

		// Start new transaction for this batch
		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to begin transaction for batch %d: %w", batchNum, err)
		}

		txQueries := r.queries.WithTx(tx)

		// Process all records in this batch
		for j, corp := range batch {
			params := db.UpsertCorporationParams{
				CorporateNumber:        corp.CorporateNumber,
				Name:                   corp.Name,
				Kana:                   r.stringToNullString(corp.Kana),
				NameEn:                 r.stringToNullString(corp.NameEn),
				SearchName:             r.stringToNullString(corp.SearchName),
				PostalCode:             r.stringToNullString(corp.PostalCode),
				Location:               r.stringToNullString(corp.Location),
				PrefectureCode:         r.stringToNullString(corp.PrefectureCode),
				Status:                 corp.Status,
				CloseDate:              r.timeToNullTime(corp.CloseDate),
				CloseCause:             r.stringToNullString(corp.CloseCause),
				RepresentativeName:     r.stringToNullString(corp.RepresentativeName),
				RepresentativePosition: r.stringToNullString(corp.RepresentativePosition),
				DateOfEstablishment:    r.timeToNullTime(corp.DateOfEstablishment),
				FoundingYear:           r.int32ToNullInt32(corp.FoundingYear),
				CapitalStock:           r.int64ToNullInt64(corp.CapitalStock),
				EmployeeNumber:         r.int32ToNullInt32(corp.EmployeeNumber),
				CompanySizeMale:        r.int32ToNullInt32(corp.CompanySizeMale),
				CompanySizeFemale:      r.int32ToNullInt32(corp.CompanySizeFemale),
				BusinessItems:          r.stringToNullString(corp.BusinessItems),
				BusinessSummary:        r.stringToNullString(corp.BusinessSummary),
				CompanyUrl:             r.stringToNullString(corp.CompanyUrl),
				QualificationGrade:     r.stringToNullString(corp.QualificationGrade),
				NumberOfActivity:       r.stringToNullString(corp.NumberOfActivity),
				UpdateDate:             r.timeToNullTime(corp.UpdateDate),
			}

			_, err := txQueries.UpsertCorporation(ctx, params)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to upsert corporation %s (batch %d, item %d): %w",
					corp.CorporateNumber, batchNum, j+1, err)
			}
		}

		// Commit this batch
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction for batch %d: %w", batchNum, err)
		}

		fmt.Printf("Committed batch %d/%d (%d corporations)\n", batchNum, totalBatches, len(batch))
	}

	fmt.Printf("Successfully upserted %d corporations in %d committed transactions\n", len(corporations), totalBatches)
	return nil
}

// convertToDomain converts a db.Corporation to domain.Corporation
func (r *CorporationRepository) convertToDomain(dbCorp *db.Corporation) *domain.Corporation {
	return &domain.Corporation{
		ID:              int64(dbCorp.ID),
		CorporateNumber: dbCorp.CorporateNumber,
		Name:            dbCorp.Name,
		Kana:            r.nullStringToString(dbCorp.Kana),
		NameEn:          r.nullStringToString(dbCorp.NameEn),
		SearchName:      r.nullStringToString(dbCorp.SearchName),
		PostalCode:      r.nullStringToString(dbCorp.PostalCode),
		Location:        r.nullStringToString(dbCorp.Location),
		PrefectureCode:  r.nullStringToString(dbCorp.PrefectureCode),
		Status:          dbCorp.Status,

		// Registration Information
		CloseDate:  r.nullTimeToTime(dbCorp.CloseDate),
		CloseCause: r.nullStringToString(dbCorp.CloseCause),

		// Representative Information
		RepresentativeName:     r.nullStringToString(dbCorp.RepresentativeName),
		RepresentativePosition: r.nullStringToString(dbCorp.RepresentativePosition),

		// Company Details
		DateOfEstablishment: r.nullTimeToTime(dbCorp.DateOfEstablishment),
		FoundingYear:        r.nullInt32ToInt32(dbCorp.FoundingYear),
		CapitalStock:        r.nullInt64ToInt64(dbCorp.CapitalStock),
		EmployeeNumber:      r.nullInt32ToInt32(dbCorp.EmployeeNumber),
		CompanySizeMale:     r.nullInt32ToInt32(dbCorp.CompanySizeMale),
		CompanySizeFemale:   r.nullInt32ToInt32(dbCorp.CompanySizeFemale),

		// Business Information
		BusinessItems:      r.nullStringToString(dbCorp.BusinessItems),
		BusinessSummary:    r.nullStringToString(dbCorp.BusinessSummary),
		CompanyUrl:         r.nullStringToString(dbCorp.CompanyUrl),
		QualificationGrade: r.nullStringToString(dbCorp.QualificationGrade),
		NumberOfActivity:   r.nullStringToString(dbCorp.NumberOfActivity),

		// gBizINFO Metadata
		UpdateDate: r.nullTimeToTime(dbCorp.UpdateDate),

		// Database Metadata
		CreatedAt: r.nullTimeToTimeValue(dbCorp.CreatedAt),
		UpdatedAt: r.nullTimeToTimeValue(dbCorp.UpdatedAt),
	}
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
