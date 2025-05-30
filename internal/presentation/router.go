package presentation

import (
	"encoding/json"
	"net/http"

	"corporation-db/internal/api"
	"corporation-db/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Router holds the HTTP router and handlers
type Router struct {
	corporationHandler *CorporationHandler
}

// NewRouter creates a new Router
func NewRouter(corporationUsecase *usecase.CorporationUsecase) *Router {
	return &Router{
		corporationHandler: NewCorporationHandler(corporationUsecase),
	}
}

// SetupRoutes configures all routes using the generated API handler
func (router *Router) SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	// Create a server interface wrapper with our handlers
	serverWrapper := &api.ServerInterfaceWrapper{
		Handler: &CombinedHandler{
			corporationHandler: router.corporationHandler,
		},
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
	}

	// Register API routes using the generated wrapper
	r.Get("/corporations", serverWrapper.GetCorporations)
	r.Get("/corporations/{corporate_number}", serverWrapper.GetCorporationsCorporateNumber)
	r.Get("/health", serverWrapper.HealthCheck)

	return r
}

// CombinedHandler implements the ServerInterface for all endpoints
type CombinedHandler struct {
	corporationHandler *CorporationHandler
}

// Corporation endpoints
func (ch *CombinedHandler) GetCorporations(w http.ResponseWriter, r *http.Request, params api.GetCorporationsParams) {
	ch.corporationHandler.GetCorporations(w, r, params)
}

func (ch *CombinedHandler) GetCorporationsCorporateNumber(w http.ResponseWriter, r *http.Request, corporateNumber string) {
	ch.corporationHandler.GetCorporationsCorporateNumber(w, r, corporateNumber)
}

// Health check endpoint
func (ch *CombinedHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
