package handler

import (
	"inventories-app/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HttpHandler struct {
	DB repository.DBRepo
}

func (hs *HttpHandler) Routes() http.Handler {

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)

	mux.Route("/v1", func(r chi.Router) {
		r.Get("/items", hs.GetItems)
		r.Get("/items/{id}", hs.GetItem)
		r.Post("/items", hs.AddItem)

		r.Put("/items/{id}", hs.UpdateItem)
		r.Delete("/items/{id}", hs.DeleteItem)
	})

	return mux
}
