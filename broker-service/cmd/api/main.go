package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = "80"

type Config struct{}

func main() {
	app := Config{}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	log.Printf("Broker service started on %s\n", port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
