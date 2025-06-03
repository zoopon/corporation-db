package domain

import (
	"context"

	"github.com/google/uuid"
)

// CorporationRepository defines the interface for corporation data access
type CorporationRepository interface {
	GetAll(ctx context.Context) ([]*Corporation, error)
	GetWithFilter(ctx context.Context, filter CorporationFilter) ([]*Corporation, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Corporation, error)
	GetByCorporateNumber(ctx context.Context, corporateNumber string) (*Corporation, error)
	Create(ctx context.Context, corp *Corporation) (*Corporation, error)
	BulkUpsert(ctx context.Context, corporations []*CreateCorporationRequest) error
}
