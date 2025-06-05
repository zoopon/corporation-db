package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Base represents a base/branch office entity
type Base struct {
	ID              uuid.UUID `json:"id" db:"id"`
	CorporationID   uuid.UUID `json:"corporation_id" db:"corporation_id"`
	CorporateNumber string    `json:"corporate_number" db:"corporate_number"`
	BaseName        *string   `json:"base_name,omitempty" db:"base_name"`
	CountryCode     string    `json:"country_code" db:"country_code"`
	PostalCode      *string   `json:"postal_code,omitempty" db:"postal_code"`
	Location        *string   `json:"location,omitempty" db:"location"`
	PhoneNumber     *string   `json:"phone_number,omitempty" db:"phone_number"`
	FaxNumber       *string   `json:"fax_number,omitempty" db:"fax_number"`
	DataObtainedAt  time.Time `json:"data_obtained_at" db:"data_obtained_at"`
	DataSourceURL   string    `json:"data_source_url" db:"data_source_url"`
	IsHeadOffice    bool      `json:"is_head_office" db:"is_head_office"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// NewBase creates a new Base with UUIDv7
func NewBase(corporationID uuid.UUID, corporateNumber string) *Base {
	return &Base{
		ID:              NewUUIDv7(),
		CorporationID:   corporationID,
		CorporateNumber: corporateNumber,
		IsHeadOffice:    false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

// NewHeadOfficeBase creates a new head office Base with UUIDv7
func NewHeadOfficeBase(corporationID uuid.UUID, corporateNumber string) *Base {
	base := NewBase(corporationID, corporateNumber)
	base.IsHeadOffice = true
	return base
}

// BaseRepository defines the interface for base persistence
type BaseRepository interface {
	Create(ctx context.Context, base *Base) (*Base, error)
	CreateBatch(ctx context.Context, bases []*Base) error
	GetByID(ctx context.Context, id uuid.UUID) (*Base, error)
	GetByCorporationID(ctx context.Context, corporationID uuid.UUID) ([]*Base, error)
	GetByCorporateNumber(ctx context.Context, corporateNumber string) ([]*Base, error)
	GetHeadOfficeByCorporateNumber(ctx context.Context, corporateNumber string) (*Base, error)
	Update(ctx context.Context, base *Base) (*Base, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByCorporationID(ctx context.Context, corporationID uuid.UUID) error
	ListAll(ctx context.Context, limit, offset int) ([]*Base, error)
	SearchByName(ctx context.Context, name string, limit, offset int) ([]*Base, error)
	GetByCountry(ctx context.Context, countryCode string, limit, offset int) ([]*Base, error)
	CountByCorporateNumber(ctx context.Context, corporateNumber string) (int64, error)
}
