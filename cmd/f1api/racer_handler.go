package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/drkkhn/f1api/pkg/f1api/model"
	"github.com/gorilla/mux"
)

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *application) createRacerHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		TeamId    string `json:"teamId"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	racer := &model.Racer{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		TeamId:    input.TeamId,
	}

	err = app.models.Racers.Insert(racer)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, racer)
}

func (app *application) getRacerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["racerId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid racer ID")
		return
	}

	racer, err := app.models.Racers.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, racer)
}

func (app *application) updateRacerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["racerId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid racer ID")
		return
	}

	racer, err := app.models.Racers.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		FirstName *string `json:"firstName"`
		LastName  *string `json:"lastName"`
		TeamId    *string `json:"teamId"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.FirstName != nil {
		racer.FirstName = *input.FirstName
	}

	if input.LastName != nil {
		racer.LastName = *input.LastName
	}

	if input.TeamId != nil {
		racer.TeamId = *input.TeamId
	}

	err = app.models.Racers.Update(racer)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, racer)
}

func (app *application) deleteRacerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["racerId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid racer ID")
		return
	}

	err = app.models.Racers.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}
