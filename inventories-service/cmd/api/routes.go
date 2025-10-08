package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)

	mux.Get("/", app.Home)

	// should return list of inventories
	mux.Get("/items", app.GetItems)
	// return invotories by id
	mux.Get("/items/{id}", nil)
	// add inventory
	mux.Post("/items", nil)
	// update invotories by id
	mux.Put("/items/{id}", nil)
	// delete inventory
	mux.Delete("/items/{id}", nil)

	return mux
}
