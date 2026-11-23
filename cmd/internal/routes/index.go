package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/zaheershaikh936/CommitLaunchBackend/cmd/internal/feature/health"
)

func ApplicationRoutes(r *chi.Mux) {
	r.Route("/health", func(r chi.Router) { health.HealthRegisterRoutes(r) })
}
