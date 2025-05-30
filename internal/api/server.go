package api

import (
	"encoding/json"
	"net/http"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

// APIServer implements the ServerInterface
type APIServer struct {
	// ここでデータベースやユースケースを注入できます
}

// NewAPIServer creates a new API server instance
func NewAPIServer() *APIServer {
	return &APIServer{}
}

// HealthCheck implements health check endpoint
func (s *APIServer) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GetUsers implements get all users endpoint
func (s *APIServer) GetUsers(w http.ResponseWriter, r *http.Request) {
	// モックデータ（後でデータベースから取得するように変更）
	now := time.Now()
	email := "john@example.com"
	name := "John Doe"
	phone := "123-456-7890"
	address := "123 Main St, City, Country"
	id := int64(1)

	users := []User{
		{
			Id:        &id,
			Name:      &name,
			Email:     (*openapi_types.Email)(&email),
			Phone:     &phone,
			Address:   &address,
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// CreateUser implements create user endpoint
func (s *APIServer) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// モックレスポンス（後でデータベースに保存するように変更）
	now := time.Now()
	id := int64(1)
	email := (*openapi_types.Email)(&req.Email)

	user := User{
		Id:        &id,
		Name:      &req.Name,
		Email:     email,
		Phone:     req.Phone,
		Address:   req.Address,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUserById implements get user by ID endpoint
func (s *APIServer) GetUserById(w http.ResponseWriter, r *http.Request, id int64) {
	// モックデータ（後でデータベースから取得するように変更）
	if id != 1 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	now := time.Now()
	email := "john@example.com"
	name := "John Doe"
	phone := "123-456-7890"
	address := "123 Main St, City, Country"

	user := User{
		Id:        &id,
		Name:      &name,
		Email:     (*openapi_types.Email)(&email),
		Phone:     &phone,
		Address:   &address,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
