package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"go-final/pkg/my-apishka/model"

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
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:Barakat2005%23@localhost/gofinal?sslmode=disable", "PostgreSQL DSN")
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

	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Обработчики маршрутов
	v1.HandleFunc("/character", app.createCharacterHandler).Methods("POST")
	v1.HandleFunc("/character/{id}", app.getCharacterHandler).Methods("GET")
	v1.HandleFunc("/character/{id}", app.updateCharacterHandler).Methods("PUT")
	v1.HandleFunc("/character/{id}", app.deleteCharacterHandler).Methods("DELETE")

	// функции по ТСИС3
	v1.HandleFunc("/charactersfilter", app.getByHouseHandler).Methods("GET")                  //по факультету
	v1.HandleFunc("/characterssorting", app.getByLastNameHandler).Methods("GET")              //по фамилиям
	v1.HandleFunc("/characterspagination", app.getCharactersPaginationHandler).Methods("GET") //устанавливается лимит на вывод данных

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
