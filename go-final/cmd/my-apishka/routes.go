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

	//для сущности коммент
	v1.HandleFunc("/comments", app.CreateCommentHandler).Methods("POST")
	v1.HandleFunc("/comments/{id}", app.GetCommentHandler).Methods("GET")
	v1.HandleFunc("/comments/{id}", app.UpdateCommentHandler).Methods("PUT")
	v1.HandleFunc("/comments/{id}", app.requirePermissions("comments:write",app.DeleteCommentHandler)).Methods("DELETE")
	
	//фильтрация,сортировка,пагинация для комментов
	v1.HandleFunc("/commentsfilter", app.getCommentsByUserIDHandler).Methods("GET")              //фильтр по айди юзера указанного в парам   
	v1.HandleFunc("/commentssorting", app.getCommentsByCharacterIDHandler).Methods("GET")        // сортинг по айди персонажей    
	v1.HandleFunc("/commentspagination", app.getCharactersPaginationHandler).Methods("GET")      //устанавливается лимит на вывод данных

	//вывод списка комментариев по айди персонажа
	v1.HandleFunc("/character/{id}/comments",app.getCharacterCommentsHandler).Methods("GET")

	//вывод списка комментариев по айди юзера
	v1.HandleFunc("/users/{id}/comments", app.getUserCommentsHandler).Methods("GET")

	return app.authenticate(r)
}