package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"orders-app/cmd/api/handler"
	"orders-app/internal/repository"
	pb "rpc"

	"google.golang.org/grpc"
)

type httpServer struct {
	addr        string
	HttpHandler handler.HttpHandler
}

func NewHttpServer(addr string, repo repository.DBRepo, conn *grpc.ClientConn) *httpServer {

	return &httpServer{
		addr: addr,
		HttpHandler: handler.HttpHandler{
			DB:               repo,
			InventoryService: pb.NewInventoryRpcClient(conn),
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
