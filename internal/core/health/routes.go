package health

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	healthhandler "github.com/vadxq/go-rest-starter/internal/core/health/handler"
)

// RegisterRoutes registers health and utility endpoints.
func RegisterRoutes(r chi.Router, handler *healthhandler.HealthHandler) {
	r.Get("/health", handler.Health)
	r.Get("/health/detailed", handler.DetailedHealth)
	r.Get("/ready", handler.Ready)
	r.Get("/live", handler.Live)

	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"version":"1.0.0"}`))
	})

	r.Route("/status", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"running"}`))
		})
	})
}
