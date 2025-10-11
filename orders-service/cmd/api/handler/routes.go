package handler

import (
	"net/http"
	"orders-app/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	pb "rpc"
)

type HttpHandler struct {
	DB               repository.DBRepo
	InventoryService pb.InventoryRpcClient
}

func (hs *HttpHandler) Routes() http.Handler {

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)

	mux.Route("/v1", func(r chi.Router) {

		r.Post("/orders", hs.AddOrder)
		r.Get("/hello", hs.Hello)

	})

	return mux
}
