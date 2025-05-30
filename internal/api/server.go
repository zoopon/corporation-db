package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

// APIServer implements the ServerInterface
type APIServer struct {
	db *sql.DB
}

// NewAPIServer creates a new API server instance
func NewAPIServer(db *sql.DB) *APIServer {
	return &APIServer{
		db: db,
	}
}

// HealthCheck implements health check endpoint
func (s *APIServer) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GetUsers implements get all users endpoint
func (s *APIServer) GetUsers(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, name, email, phone, address, created_at, updated_at FROM users ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(r.Context(), query)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var phone, address sql.NullString
		var createdAt, updatedAt sql.NullTime
		var id int64
		var name, email string

		err := rows.Scan(&id, &name, &email, &phone, &address, &createdAt, &updatedAt)
		if err != nil {
			http.Error(w, "Failed to scan user", http.StatusInternalServerError)
			return
		}

		// Convert to API types
		user.Id = &id
		user.Name = &name
		emailType := openapi_types.Email(email)
		user.Email = &emailType

		if phone.Valid {
			user.Phone = &phone.String
		}
		if address.Valid {
			user.Address = &address.String
		}
		if createdAt.Valid {
			user.CreatedAt = &createdAt.Time
		}
		if updatedAt.Valid {
			user.UpdatedAt = &updatedAt.Time
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Failed to iterate users", http.StatusInternalServerError)
		return
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

	// Prepare nullable fields
	var phone, address sql.NullString
	if req.Phone != nil {
		phone = sql.NullString{String: *req.Phone, Valid: true}
	}
	if req.Address != nil {
		address = sql.NullString{String: *req.Address, Valid: true}
	}

	query := `INSERT INTO users (name, email, phone, address) VALUES ($1, $2, $3, $4) 
			  RETURNING id, name, email, phone, address, created_at, updated_at`

	var user User
	var id int64
	var name, email string
	var dbPhone, dbAddress sql.NullString
	var createdAt, updatedAt sql.NullTime

	err := s.db.QueryRowContext(r.Context(), query, req.Name, req.Email, phone, address).
		Scan(&id, &name, &email, &dbPhone, &dbAddress, &createdAt, &updatedAt)

	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Convert to API types
	user.Id = &id
	user.Name = &name
	emailType := openapi_types.Email(email)
	user.Email = &emailType

	if dbPhone.Valid {
		user.Phone = &dbPhone.String
	}
	if dbAddress.Valid {
		user.Address = &dbAddress.String
	}
	if createdAt.Valid {
		user.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		user.UpdatedAt = &updatedAt.Time
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUserById implements get user by ID endpoint
func (s *APIServer) GetUserById(w http.ResponseWriter, r *http.Request, id int64) {
	query := `SELECT id, name, email, phone, address, created_at, updated_at FROM users WHERE id = $1 LIMIT 1`

	var user User
	var dbId int64
	var name, email string
	var phone, address sql.NullString
	var createdAt, updatedAt sql.NullTime

	err := s.db.QueryRowContext(r.Context(), query, id).
		Scan(&dbId, &name, &email, &phone, &address, &createdAt, &updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	// Convert to API types
	user.Id = &dbId
	user.Name = &name
	emailType := openapi_types.Email(email)
	user.Email = &emailType

	if phone.Valid {
		user.Phone = &phone.String
	}
	if address.Valid {
		user.Address = &address.String
	}
	if createdAt.Valid {
		user.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		user.UpdatedAt = &updatedAt.Time
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
