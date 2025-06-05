package presentation

import (
	"encoding/json"
	"net/http"

	"corporation-db/internal/api"
	"corporation-db/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Router holds the HTTP router and handlers
type Router struct {
	corporationHandler *CorporationHandler
}

// NewRouter creates a new Router
func NewRouter(corporationUsecase *usecase.CorporationUsecase, baseUsecase *usecase.BaseUsecase, fetchBasesUsecase *usecase.FetchBasesUseCase) *Router {
	return &Router{
		corporationHandler: NewCorporationHandler(corporationUsecase, baseUsecase, fetchBasesUsecase),
	}
}

// SetupRoutes configures all routes using the generated API handler
func (router *Router) SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	// CORS middleware - Allow cross-origin requests
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},                      // Allow all origins for development
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"}, // Allow POST for fetch-bases
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	// Create a combined handler that implements ServerInterface
	combinedHandler := &CombinedHandler{
		corporationHandler: router.corporationHandler,
	}

	// Use the generated API handler with options
	apiHandler := api.HandlerWithOptions(combinedHandler, api.ChiServerOptions{
		BaseURL:    "",
		BaseRouter: r,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
	})

	return apiHandler.(*chi.Mux)
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

// FetchCorporationBases implements the POST /corporations/{corporate_number}/fetch-bases endpoint
func (ch *CombinedHandler) FetchCorporationBases(w http.ResponseWriter, r *http.Request, corporateNumber string) {
	ch.corporationHandler.FetchCorporationBases(w, r, corporateNumber)
}

// Health check endpoint
func (ch *CombinedHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
