package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)

	mux.Route("/v1", func(r chi.Router) {
		r.Get("/items", app.GetItems)
		r.Get("/items/{id}", app.GetItem)
		r.Post("/items", app.AddItem)

		r.Put("/items/{id}", nil)
		r.Delete("/items/{id}", app.DeleteItem)
	})

	return mux
}
