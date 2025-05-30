package domain

import (
	"context"
)

// CorporationRepository defines the interface for corporation data access
type CorporationRepository interface {
	GetAll(ctx context.Context) ([]*Corporation, error)
	GetWithFilter(ctx context.Context, filter CorporationFilter) ([]*Corporation, int64, error)
	GetByID(ctx context.Context, id int64) (*Corporation, error)
	GetByCorporateNumber(ctx context.Context, corporateNumber string) (*Corporation, error)
	BulkUpsert(ctx context.Context, corporations []*CreateCorporationRequest) error
}
