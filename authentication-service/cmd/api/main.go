package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const port = "80"

var counts = 0

type Config struct {
	Repo data.Repository
}

func main() {
	conn := connectDB()
	if conn == nil {
		log.Panic("can not connect to postgres!")
	}

	app := Config{
		Repo: data.NewPostgresRepo(conn),
	}

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

func connectDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		db, err := openDB(dsn)
		if err != nil {
			log.Println("postgres not yet ready...", err)
			counts++
		} else {
			return db
		}

		if counts > 10 {
			log.Println("counts > 10, too many fails")
			return nil
		}

		log.Println("retry connecting postgres in 2 seconds...")
		time.Sleep(2 * time.Second)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
