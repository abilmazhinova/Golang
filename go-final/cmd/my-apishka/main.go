package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"sync"

	"go-final/pkg/my-apishka/model"
	// "github.com/codev0/inft3212-6/pkg/abr-plus/model/filler"
	"go-final/pkg/jsonlog"
	"go-final/pkg/vcs"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/peterbourgon/ff/v3"

	_ "github.com/lib/pq"
)

// Set version of application corresponding to value of vcs.Version.
var (
	version = vcs.Version()
)

type config struct {
	port       int
	env        string
	fill       bool
	migrations string
	db         struct {
		dsn string
	}
}

type application struct {
	config config
	models model.Models
	logger *jsonlog.Logger
	wg     sync.WaitGroup
}

func main() {
	fs := flag.NewFlagSet("demo-app", flag.ContinueOnError)

	var (
		cfg        config
		fill       = fs.Bool("fill", false, "Fill database with dummy data")
		migrations = fs.String("migrations", "", "Path to migration files folder. If not provided, migrations do not applied")
		port       = fs.Int("port", 8081, "API server port")
		env        = fs.String("env", "development", "Environment (development|staging|production)")
		dbDsn      = fs.String("dsn", "postgres://postgres:Barakat2005%23@localhost/gofinal?sslmode=disable", "PostgreSQL DSN")
	)

	// Init logger
	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVars()); err != nil {
		logger.PrintFatal(err, nil)
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}

	cfg.port = *port
	cfg.env = *env
	cfg.fill = *fill
	cfg.db.dsn = *dbDsn
	cfg.migrations = *migrations

	logger.PrintInfo("starting application with configuration", map[string]string{
		"port":       fmt.Sprintf("%d", cfg.port),
		"fill":       fmt.Sprintf("%t", cfg.fill),
		"env":        cfg.env,
		"db":         cfg.db.dsn,
		"migrations": cfg.migrations,
	})

	// Connect to DB
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintError(err, nil)
		return
	}
	// Defer a call to db.Close() so that the connection pool is closed before the main()
	// function exits.
	defer func() {
		if err := db.Close(); err != nil {
			logger.PrintFatal(err, nil)
		}
	}()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
		logger: logger,
	}

	// if cfg.fill {
	// 	err = filler.PopulateDatabase(app.models)
	// 	if err != nil {
	// 		logger.PrintFatal(err, nil)
	// 		return
	// 	}
	// }

	// Call app.server() to start the server.
	if err := app.serve(); err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config // struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// https://github.com/golang-migrate/migrate?tab=readme-ov-file#use-in-your-go-project
	if cfg.migrations != "" {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return nil, err
		}
		m, err := migrate.NewWithDatabaseInstance(
			cfg.migrations,
			"postgres", driver)
		if err != nil {
			return nil, err
		}
		m.Up()
	}

	return db, nil
}






// package main

// import (
// 	"database/sql"
// 	"flag"
// 	"log"
// 	"net/http"

// 	"go-final/pkg/my-apishka/model"

// 	"github.com/gorilla/mux"

// 	_ "github.com/lib/pq"
// )

// type config struct {
// 	port string
// 	env  string
// 	db   struct {
// 		dsn string
// 	}
// }

// type application struct {
// 	config config
// 	models model.Models
// }

// func main() {
// 	var cfg config
// 	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
// 	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
// 	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:Barakat2005%23@localhost/gofinal?sslmode=disable", "PostgreSQL DSN")
// 	flag.Parse()

// 	db, err := openDB(cfg)
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}
// 	defer db.Close()

// 	app := &application{
// 		config: cfg,
// 		models: model.NewModels(db),
// 	}

// 	app.run()
// }

// func (app *application) run() {
// 	r := mux.NewRouter()

// 	v1 := r.PathPrefix("/api/v1").Subrouter()

// 	// Обработчики маршрутов
// 	v1.HandleFunc("/character", app.createCharacterHandler).Methods("POST")
// 	v1.HandleFunc("/character/{id}", app.getCharacterHandler).Methods("GET")
// 	v1.HandleFunc("/character/{id}", app.updateCharacterHandler).Methods("PUT")
// 	v1.HandleFunc("/character/{id}", app.deleteCharacterHandler).Methods("DELETE")

// 	// функции по ТСИС3
// 	v1.HandleFunc("/charactersfilter", app.getByHouseHandler).Methods("GET")                  //по факультету
// 	v1.HandleFunc("/characterssorting", app.getByLastNameHandler).Methods("GET")              //по фамилиям
// 	v1.HandleFunc("/characterspagination", app.getCharactersPaginationHandler).Methods("GET") //устанавливается лимит на вывод данных

// 	//для сущности юзера
// 	v1.HandleFunc("/users",app.registerUserHandler).Methods("POST")

// 	log.Printf("Starting server on %s\n", app.config.port)
// 	err := http.ListenAndServe(app.config.port, r)
// 	log.Fatal(err)
// }

// func openDB(cfg config) (*sql.DB, error) {
// 	db, err := sql.Open("postgres", cfg.db.dsn)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = db.Ping()
// 	if err != nil {
// 		db.Close()
// 		return nil, err
// 	}

// 	return db, nil
// }
