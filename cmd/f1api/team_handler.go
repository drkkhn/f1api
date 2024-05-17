package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/drkkhn/f1api/pkg/f1api/model"
	"github.com/drkkhn/f1api/pkg/f1api/validator"
	"github.com/gorilla/mux"
)

func (app *application) createTeamHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
		Car  string `json:"car"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	team := &model.Team{
		Name: input.Name,
		Car:  input.Car,
	}

	v := validator.New()

	if model.ValidateTeam(v, team); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Teams.Insert(team)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, team)
}

func (app *application) getTeamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["teamId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	team, err := app.models.Teams.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, team)
}

func (app *application) getTeamRacer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["teamId"]

	fmt.Println("TEAM RACERS")

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	racers, err := app.models.Teams.GetTeamRacers(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	if len(*racers) >= 1 {
		app.respondWithJSON(w, http.StatusOK, racers)
	} else {
		app.respondWithJSON(w, http.StatusNotFound, "there is no racers in the team")
	}
}

func (app *application) updateTeamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["teamId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	team, err := app.models.Teams.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Name *string `json:"name"`
		Car  *string `json:"car"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Name != nil {
		team.Name = *input.Name
	}

	if input.Car != nil {
		team.Car = *input.Car
	}

	err = app.models.Teams.Update(team)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, team)
}

func (app *application) deleteTeamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["teamId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	err = app.models.Teams.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
