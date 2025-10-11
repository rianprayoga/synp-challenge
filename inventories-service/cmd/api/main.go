package main

import (
	"flag"
	"inventories-app/cmd/api/httpserver"
	"inventories-app/cmd/api/rpcserver"
	"inventories-app/internal/repository"
	"inventories-app/internal/repository/dbrepo"
	"log"
)

type application struct {
	DSN      string
	HttpPort string
	GrpcPort string
	DB       repository.DBRepo
}

func main() {
	var app application

	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=inventories timezone=UTC", "Postgres connection string")
	flag.StringVar(&app.HttpPort, "http-port", "8081", "Port for http inventories service")
	flag.StringVar(&app.GrpcPort, "grpc-port", "8082", "Port for grpc inventories service")
	flag.Parse()

	conn, err := app.connectDb()
	if err != nil {
		log.Fatal(err)

	}
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer conn.Close()

	gRPCServer := rpcserver.NewGRPCServer(app.GrpcPort, app.DB)
	go gRPCServer.Run()

	httpServer := httpserver.NewHttpServer(app.HttpPort, app.DB)
	httpServer.Run()

}
