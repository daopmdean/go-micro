package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

const port = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	app := Config{}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	log.Printf("Authentication service started on %s\n", port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
