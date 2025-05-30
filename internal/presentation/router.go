package presentation

import (
	"encoding/json"
	"net/http"

	"corporation-db/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Router holds the HTTP router and handlers
type Router struct {
	userHandler *UserHandler
}

// NewRouter creates a new Router
func NewRouter(userUsecase *usecase.UserUsecase) *Router {
	return &Router{
		userHandler: NewUserHandler(userUsecase),
	}
}

// SetupRoutes configures all routes
func (router *Router) SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	// Health check
	r.Get("/health", router.healthCheck)

	// User routes
	r.Route("/users", func(r chi.Router) {
		r.Get("/", router.userHandler.GetUsers)
		r.Post("/", router.userHandler.CreateUser)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", router.userHandler.GetUserByID)
			r.Put("/", router.userHandler.UpdateUser)
			r.Delete("/", router.userHandler.DeleteUser)
		})
	})

	return r
}

// healthCheck handles health check endpoint
func (router *Router) healthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
