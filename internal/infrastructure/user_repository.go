package infrastructure

import (
	"context"
	"database/sql"
	"fmt"

	"corporation-db/internal/domain"
	"corporation-db/internal/infrastructure/db"

	_ "github.com/lib/pq"
)

// UserRepository implements domain.UserRepository
type UserRepository struct {
	db      *sql.DB
	queries *db.Queries
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(database *sql.DB) domain.UserRepository {
	return &UserRepository{
		db:      database,
		queries: db.New(database),
	}
}

// GetAll retrieves all users
func (r *UserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	dbUsers, err := r.queries.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	users := make([]*domain.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = r.convertToUser(dbUser)
	}

	return users, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	dbUser, err := r.queries.GetUserByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return r.convertToUser(dbUser), nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return r.convertToUser(dbUser), nil
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	dbUser, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return r.convertToUser(dbUser), nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, id int64, req *domain.UpdateUserRequest) (*domain.User, error) {
	dbUser, err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:    int32(id),
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return r.convertToUser(dbUser), nil
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	err := r.queries.DeleteUser(ctx, int32(id))
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// convertToUser converts database model to domain model
func (r *UserRepository) convertToUser(dbUser db.User) *domain.User {
	user := &domain.User{
		ID:    int64(dbUser.ID),
		Name:  dbUser.Name,
		Email: dbUser.Email,
	}

	if dbUser.CreatedAt.Valid {
		user.CreatedAt = dbUser.CreatedAt.Time
	}
	if dbUser.UpdatedAt.Valid {
		user.UpdatedAt = dbUser.UpdatedAt.Time
	}

	return user
}
