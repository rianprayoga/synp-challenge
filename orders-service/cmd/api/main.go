package main

import (
	"flag"
	"log"
	"orders-app/cmd/api/httpserver"
	"orders-app/internal/repository"
	"orders-app/internal/repository/dbrepo"
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

	httpServer := httpserver.NewHttpServer(app.HttpPort, app.DB)
	httpServer.Run()
}
