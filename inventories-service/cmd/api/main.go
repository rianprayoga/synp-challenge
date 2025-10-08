package main

import (
	"flag"
	"fmt"
	"inventories-app/internal/repository"
	"inventories-app/internal/repository/dbrepo"
	"log"
	"net/http"
)

type application struct {
	DSN  string
	Port string
	DB   repository.DBRepo
}

func main() {
	var app application

	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=inventories timezone=UTC", "Postgres connection string")
	flag.StringVar(&app.Port, "port", "8081", "Port for inventories service")
	flag.Parse()

	log.Println("starting app on port ", app.Port)

	conn, err := app.connectDb()
	if err != nil {
		log.Fatal(err)

	}
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer conn.Close()

	err = http.ListenAndServe(fmt.Sprintf(":%s", app.Port), app.routes())
	if err != nil {
		log.Fatal(err)
	}

}
