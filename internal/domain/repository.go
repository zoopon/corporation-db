package domain

import (
	"context"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	GetAll(ctx context.Context) ([]*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, req *CreateUserRequest) (*User, error)
	Update(ctx context.Context, id int64, req *UpdateUserRequest) (*User, error)
	Delete(ctx context.Context, id int64) error
}

// CorporationRepository defines the interface for corporation data access
type CorporationRepository interface {
	GetAll(ctx context.Context) ([]*Corporation, error)
	GetWithFilter(ctx context.Context, filter CorporationFilter) ([]*Corporation, int64, error)
	GetByID(ctx context.Context, id int64) (*Corporation, error)
	GetByCorporateNumber(ctx context.Context, corporateNumber string) (*Corporation, error)
	BulkUpsert(ctx context.Context, corporations []*CreateCorporationRequest) error
}
