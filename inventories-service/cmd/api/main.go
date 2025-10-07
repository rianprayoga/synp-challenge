package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

type application struct {
	DSN  string
	Port string
}

func main() {
	var app application

	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=inventories timezone=UTC", "Postgres connection string")
	flag.StringVar(&app.Port, "port", "8081", "Port for inventories service")
	flag.Parse()

	err := http.ListenAndServe(fmt.Sprintf(":%s", app.Port), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("starting app on port ", app.Port)

}
