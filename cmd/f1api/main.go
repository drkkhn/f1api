package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/drkkhn/f1api/pkg/f1api/model"
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
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:0000@localhost/f1?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	// Connect to DB
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

	// Racer Singleton
	// Create a new Racer
	v1.HandleFunc("/racers", app.createRacerHandler).Methods("POST")
	// Get a specific Racer
	v1.HandleFunc("/racers/{racerId:[0-9]+}", app.getRacerHandler).Methods("GET")
	// Update a specific Racer
	v1.HandleFunc("/racers/{racerId:[0-9]+}", app.updateRacerHandler).Methods("PUT")
	// Delete a specific Racer
	v1.HandleFunc("/racers/{racerId:[0-9]+}", app.deleteRacerHandler).Methods("DELETE")

	// Create a new racer
	v1.HandleFunc("/teams", app.createTeamHandler).Methods("POST")
	// Get a specific racer
	v1.HandleFunc("/teams/{teamId:[0-9]+}", app.getTeamHandler).Methods("GET")
	// Update a specific racer
	v1.HandleFunc("/teams/{teamId:[0-9]+}", app.updateTeamHandler).Methods("PUT")
	// Delete a specific racer
	v1.HandleFunc("/teams/{teamId:[0-9]+}", app.deleteTeamHandler).Methods("DELETE")

	log.Printf("Starting server on %s\n", app.config.port)
	err := http.ListenAndServe(app.config.port, r)
	log.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config // struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
