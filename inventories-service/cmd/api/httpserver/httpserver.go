package httpserver

import (
	"fmt"
	"inventories-app/cmd/api/handler"
	"inventories-app/internal/repository"
	"log"
	"net/http"
)

type httpServer struct {
	addr        string
	HttpHandler handler.HttpHandler
}

func NewHttpServer(addr string, repo repository.DBRepo) *httpServer {

	return &httpServer{
		addr: addr,
		HttpHandler: handler.HttpHandler{
			DB: repo,
		},
	}
}

func (hs *httpServer) Run() error {

	err := http.ListenAndServe(fmt.Sprintf(":%s", hs.addr), hs.HttpHandler.Routes())
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
