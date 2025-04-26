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

const webPort = "80"

var dbConnectionAttempts int64

type Config struct {
	Repo data.Repository
	Client *http.Client
}

func main() {
	log.Println("Starting authentication service")

	// Connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	app := Config{
		Client: &http.Client{},
	}

	serve := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := serve.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func openDB(dsn string) (* sql.DB, error) {
	db, error := sql.Open("pgx", dsn)
	if error != nil {
		return nil, error
	}

	error = db.Ping()
	if error != nil {
		return nil, error
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, error := openDB(dsn)
		if error != nil {
			log.Println("Postgres not yet ready...")
			dbConnectionAttempts++
		} else {
			log.Println("Connected to Postgres!")
			// reset dbconnectionAttempts?
			// dbConnectionAttempts = 0
			return connection
		}

		if dbConnectionAttempts > 10 {
			log.Println(error)
			return nil
		}

		log.Println("Backing off for 2 seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}

func (app *Config) setupRepo(conn *sql.DB) {
	db := data.NewPostgresRepository(conn)
	app.Repo = db	
}