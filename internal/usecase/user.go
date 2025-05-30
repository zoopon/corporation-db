package usecase

import (
	"context"
	"fmt"

	"corporation-db/internal/domain"
)

// UserUsecase handles user business logic
type UserUsecase struct {
	userRepo domain.UserRepository
}

// NewUserUsecase creates a new UserUsecase
func NewUserUsecase(userRepo domain.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

// GetUsers retrieves all users
func (u *UserUsecase) GetUsers(ctx context.Context) ([]*domain.User, error) {
	return u.userRepo.GetAll(ctx)
}

// GetUserByID retrieves a user by ID
func (u *UserUsecase) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return user, nil
}

// CreateUser creates a new user
func (u *UserUsecase) CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	// Check if user with this email already exists
	existingUser, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	user, err := u.userRepo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

// UpdateUser updates an existing user
func (u *UserUsecase) UpdateUser(ctx context.Context, id int64, req *domain.UpdateUserRequest) (*domain.User, error) {
	// Check if user exists
	_, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	user, err := u.userRepo.Update(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return user, nil
}

// DeleteUser deletes a user
func (u *UserUsecase) DeleteUser(ctx context.Context, id int64) error {
	// Check if user exists
	_, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	err = u.userRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
