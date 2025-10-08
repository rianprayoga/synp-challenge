package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)

	mux.Get("/items", app.GetItems)
	mux.Get("/items/{id}", app.GetItem)
	// add inventory
	mux.Post("/items", nil)
	// update invotories by id
	mux.Put("/items/{id}", nil)
	// delete inventory
	mux.Delete("/items/{id}", nil)

	return mux
}
