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

	// Обработчики маршрутов
	v1.HandleFunc("/character", app.createCharacterHandler).Methods("POST")
	v1.HandleFunc("/character/{id}", app.getCharacterHandler).Methods("GET")
	v1.HandleFunc("/character/{id}", app.updateCharacterHandler).Methods("PUT")
	//для специальных пользователей
	v1.HandleFunc("/characters/{id}", app.requirePermissions("characters:write",app.deleteCharacterHandler)).Methods("DELETE")
	// v1.HandleFunc("/character/{id}", app.deleteCharacterHandler).Methods("DELETE")

	// функции по ТСИС3
	v1.HandleFunc("/charactersfilter", app.getByHouseHandler).Methods("GET")                  //по факультету
	v1.HandleFunc("/characterssorting", app.getByLastNameHandler).Methods("GET")              //по фамилиям
	v1.HandleFunc("/characterspagination", app.getCharactersPaginationHandler).Methods("GET") //устанавливается лимит на вывод данных

	//для сущности юзера
	v1.HandleFunc("/users",app.registerUserHandler).Methods("POST")
	v1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	v1.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")


	return app.authenticate(r)
}