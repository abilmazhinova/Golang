package main

import (
	"encoding/json"
	"github.com/abilmazhinova/Golang/final-project/pkg/my-apishka/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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

func (app *application) createCharacterHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName    string `json:"firstName"`
		LastName     string `json:"lastName"`
		House        string `json:"house"`
		OriginStatus string `json:"originStatus"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}

	character := &model.Character{
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		House:        input.House,
		OriginStatus: input.OriginStatus,
	}

	err = app.models.Questionnaires.Insert(character)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, character)
}

func (app *application) getCHaracterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["ID"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid questionnaire ID")
		return
	}

	character, err := app.models.Questionnaires.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, character)
}

func (app *application) updateCharacterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["iD"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid questionnaire ID")
		return
	}

	character, err := app.models.Characters.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		FirstName *string `json:"firstName"`
		LastName  *string `json:"lastName"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.FirstName != nil {
		character.Topic = *input.FirstName
	}

	if input.LastName != nil {
		character.Questions = *input.LastName
	}

	err = app.models.Questionnaires.Update(character)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
	}

	app.respondWithJSON(w, http.StatusOK, character)
}

func (app *application) deleteCharacterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["ID"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid questionnaire ID")
		return
	}

	err = app.models.Questionnaires.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
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