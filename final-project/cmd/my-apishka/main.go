package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/abilmazhinova/Golang/final-project/pkg/my-apishka/model"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models model.Models
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:delete@localhost/go-final?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
	}

	app.run()
}

func (app *application) run() {
    r := mux.NewRouter()

    // Обработчики маршрутов
    r.HandleFunc("/users", app.createCharacterHandler).Methods("POST")
    r.HandleFunc("/users/{id}", app.getCHaracterHandler).Methods("GET")
    r.HandleFunc("/users/{id}", app.updateCharacterHandler).Methods("PUT")
    r.HandleFunc("/users/{id}", app.deleteCharacterHandler).Methods("DELETE")

    log.Printf("Starting server on %s\n", app.config.port)
    err := http.ListenAndServe(app.config.port, r)
    log.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
    db, err := sql.Open("postgres", cfg.db.dsn)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        db.Close()
        return nil, err
    }

    return db, nil
}
