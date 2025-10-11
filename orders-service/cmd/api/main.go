package main

import (
	"flag"
	"log"
	"orders-app/cmd/api/httpserver"
	"orders-app/internal/repository"
	"orders-app/internal/repository/dbrepo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type application struct {
	DSN               string
	HttpPort          string
	InventoryGrpcAddr string
	DB                repository.DBRepo
}

func main() {
	var app application

	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=orders timezone=UTC", "Postgres connection string")
	flag.StringVar(&app.HttpPort, "port", "8083", "Port for orders service")
	flag.StringVar(&app.InventoryGrpcAddr, "inventory-grpc", "localhost:8082", "Address for inventory grpc")
	flag.Parse()

	conn, err := app.connectDb()
	if err != nil {
		log.Fatal(err)

	}
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer conn.Close()

	grpcConn, err := grpc.NewClient(app.InventoryGrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	httpServer := httpserver.NewHttpServer(app.HttpPort, app.DB, grpcConn)
	httpServer.Run()
}
