package main

import (
	"encoding/json"
	"go-final/pkg/my-apishka/model"
	// "go-final/pkg/my-apishka/validator"
	"net/http"
	"strconv"

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

func (app *application) createCharacterHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName    string `json:"FirstName"`
		LastName     string `json:"LastName"`
		House        string `json:"House"`
		OriginStatus string `json:"OriginStatus"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	character := &model.Character{
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		House:        input.House,
		OriginStatus: input.OriginStatus,
	}

	err = app.models.Characters.Insert(character)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, character)
}

func (app *application) getCharacterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid character ID")
		return
	}

	character, err := app.models.Characters.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, character)
}

func (app *application) updateCharacterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid character ID")
		return
	}

	character, err := app.models.Characters.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		FirstName *string `json:"FirstName"`
		LastName  *string `json:"LastName"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.FirstName != nil {
		character.FirstName = *input.FirstName
	}

	if input.LastName != nil {
		character.LastName = *input.LastName
	}

	err = app.models.Characters.Update(character)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
	}

	app.respondWithJSON(w, http.StatusOK, character)
}

func (app *application) deleteCharacterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid character ID")
		return
	}
 
	err = app.models.Characters.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}


func (app *application) getByHouseHandler(w http.ResponseWriter, r *http.Request) {
	house := r.URL.Query().Get("house")

	if house == "" {
		app.respondWithError(w, http.StatusBadRequest, "Please, write house name and try again.")
		return
	}

	characters, err := app.models.Characters.GetByHouse(house)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Fail, please try again.")
		return
	}

	app.respondWithJSON(w, http.StatusOK, characters)
}

func (app *application) getByLastNameHandler(w http.ResponseWriter, r *http.Request) {
	characters, err := app.models.Characters.GetByLastName()
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Fail, please try again.")
		return
	}

	app.respondWithJSON(w, http.StatusOK, characters)
}

func (app *application) getCharactersPaginationHandler(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Wrong limit parameter.")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Wrong offset parameter.")
		return
	}

	characters, err := app.models.Characters.GetCharactersPagination(limit, offset)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Fail, please try again.")
		return
	}

	app.respondWithJSON(w, http.StatusOK, characters)
}
