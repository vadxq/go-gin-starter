package user

import (
	"github.com/go-chi/chi/v5"

	userhandler "github.com/vadxq/go-rest-starter/internal/core/user/handler"
	custommiddleware "github.com/vadxq/go-rest-starter/internal/transport/middleware"
)

// RegisterRoutes registers user endpoints under the current group.
func RegisterRoutes(r chi.Router, handler *userhandler.UserHandler) {
	r.Route("/users", func(r chi.Router) {
		// Collection
		r.Get("/", handler.ListUsers)
		r.With(custommiddleware.RequireRole("admin")).Post("/", handler.CreateUser)

		// Resource
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetUser)
			r.Put("/", handler.UpdateUser)
			r.Delete("/", handler.DeleteUser)
		})
	})
}
