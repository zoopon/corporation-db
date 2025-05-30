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
