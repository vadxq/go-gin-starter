package auth

import (
	"github.com/go-chi/chi/v5"

	authhandler "github.com/vadxq/go-rest-starter/internal/core/auth/handler"
)

// RegisterPublicRoutes registers public auth endpoints under the current group.
func RegisterPublicRoutes(r chi.Router, handler *authhandler.AuthHandler) {
	// Auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", handler.Login)
		r.Post("/refresh", handler.RefreshToken)
	})
}

// RegisterProtectedRoutes registers protected auth endpoints under the current group.
func RegisterProtectedRoutes(r chi.Router, handler *authhandler.AuthHandler) {
	r.Route("/account", func(r chi.Router) {
		r.Post("/logout", handler.Logout)
	})
}
