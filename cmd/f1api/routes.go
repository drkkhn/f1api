package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// routes is our main application's router.
func (app *application) routes() http.Handler {
	r := mux.NewRouter()
	// Convert the app.notFoundResponse helper to a http.Handler using the http.HandlerFunc()
	// adapter, and then set it as the custom error handler for 404 Not Found responses.
	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	// Convert app.methodNotAllowedResponse helper to a http.Handler and set it as the custom
	// error handler for 405 Method Not Allowed responses
	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	// healthcheck
	r.HandleFunc("/api/v1/healthcheck", app.healthcheckHandler).Methods("GET")

	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Users
	v1.HandleFunc("/racers", app.requirePermission("racers:write", app.createRacerHandler)).Methods("POST")
	// Get a list of Racers
	v1.HandleFunc("/racers", app.requirePermission("racers:read", app.getAllRacersHandler)).Methods("GET")
	// Get a specific Racer
	v1.HandleFunc("/racers/{racerId:[0-9]+}", app.requirePermission("racers:read", app.getRacerHandler)).Methods("GET")
	// Update a specific Racer
	v1.HandleFunc("/racers/{racerId:[0-9]+}", app.requirePermission("racers:write", app.updateRacerHandler)).Methods("PUT")
	// Delete a specific Racer
	v1.HandleFunc("/racers/{racerId:[0-9]+}", app.requirePermission("racers:write", app.deleteRacerHandler)).Methods("DELETE")

	// Create a new racer
	v1.HandleFunc("/teams", app.requirePermission("racers:write", app.createTeamHandler)).Methods("POST")
	// Get racer with teamID
	v1.HandleFunc("/teams/{teamId:[0-9]+}/racers", app.requirePermission("racers:read", app.getTeamRacer)).Methods("GET")
	// Get a specific racer
	v1.HandleFunc("/teams/{teamId:[0-9]+}", app.requirePermission("racers:read", app.getTeamHandler)).Methods("GET")
	// Update a specific racer
	v1.HandleFunc("/teams/{teamId:[0-9]+}", app.requirePermission("racers:write", app.updateTeamHandler)).Methods("PUT")
	// Delete a specific racer
	v1.HandleFunc("/teams/{teamId:[0-9]+}", app.requirePermission("racers:write", app.deleteTeamHandler)).Methods("DELETE")

	// Members
	v1.HandleFunc("/members", app.registerMemberHandler).Methods("POST")
	v1.HandleFunc("/members/activated", app.activateMemberHandler).Methods("PUT")
	v1.HandleFunc("/tokens/authentication", app.createAuthenticationTokenHandler).Methods("POST")

	// Wrap the router with the panic recovery middleware and rate limit middleware.}
	return app.authenticate(r)
}
