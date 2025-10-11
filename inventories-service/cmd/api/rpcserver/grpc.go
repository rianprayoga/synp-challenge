package rpcserver

import (
	"fmt"
	"inventories-app/cmd/api/handler"
	"inventories-app/internal/repository"
	"log"
	"net"

	"google.golang.org/grpc"
)

type gRPCServer struct {
	addr string
	DB   repository.DBRepo
}

func NewGRPCServer(addr string, repo repository.DBRepo) *gRPCServer {
	return &gRPCServer{addr: addr, DB: repo}
}

func (gs *gRPCServer) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", gs.addr))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	handler.NewGrpcItemService(s, gs.DB)

	return s.Serve(l)
}
