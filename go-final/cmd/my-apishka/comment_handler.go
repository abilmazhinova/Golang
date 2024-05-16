package main

import (
	// "encoding/json"
	// "database/sql"
	"net/http"
	"strconv"

	"go-final/pkg/my-apishka/model"
	// "go-final/pkg/my-apishka/validator"

	"github.com/gorilla/mux"
)

// CommentHandler хранит методы для обработки HTTP запросов, связанных с комментариями.
type CommentHandler struct {
	Model *model.CommentModel
}

// CreateCommentHandler обрабатывает запрос на создание нового комментария.
func (app *application) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UsernameID  int64  `json:"UsernameID"`
		Comment     string `json:"Comment"`
		CharacterID int64  `json:"CharacterID"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	comment := &model.Comment{
		UsernameID:  input.UsernameID,
		Comment:     input.Comment,
		CharacterID: input.CharacterID,
	}

	err = app.models.Comments.CreateComment(comment)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	// comm := r.FormValue("Comment")

	// // Validate comment using the NotEmpty function
	// if !validator.NotEmpty(comm) {
	//     app.respondWithError(w, http.StatusBadRequest, "Comment must not be empty")
	//     return
	// }

	app.respondWithJSON(w, http.StatusCreated, comment)
}

// GetCommentByIDHandler обрабатывает запрос на получение комментария по его ID.
func (app *application) GetCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	comment, err := app.models.Comments.GetCommentByID(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, comment)
}

// UpdateCommentHandler обрабатывает запрос на обновление существующего комментария.
func (app *application) UpdateCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	comment, err := app.models.Comments.GetCommentByID(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Comment *string `json:"Comment"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Comment != nil {
		comment.Comment = *input.Comment
	}

	err = app.models.Comments.UpdateComment(comment)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
	}

	app.respondWithJSON(w, http.StatusOK, comment)
}

// DeleteCommentHandler обрабатывает запрос на удаление комментария по его ID.
func (app *application) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	err = app.models.Comments.DeleteCommentByID(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// фильтр,сорт,пагинация
func (app *application) getCommentsByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userID")

	if userIDStr == "" {
		app.respondWithError(w, http.StatusBadRequest, "Please provide a valid userID.")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid userID format.")
		return
	}

	comments, err := app.models.Comments.GetCommentsByUserID(userID)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to fetch comments.")
		return
	}

	app.respondWithJSON(w, http.StatusOK, comments)
}

func (app *application) getCommentsByCharacterIDHandler(w http.ResponseWriter, r *http.Request) {
	// Здесь предполагается, что фильтрация по айди персонажа = айди персонажа
	comments, err := app.models.Comments.GetCommentsByCharacter()
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to fetch comments.")
		return
	}

	app.respondWithJSON(w, http.StatusOK, comments)
}

func (app *application) getCommentsWithPaginationHandler(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid limit parameter.")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid offset parameter.")
		return
	}

	comments, err := app.models.Comments.GetCommentsPagination(limit, offset)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to fetch comments.")
		return
	}

	app.respondWithJSON(w, http.StatusOK, comments)
}

// выводим список комментариев по айди персонажа
func (app *application) getCharacterCommentsHandler(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	characterID, err := strconv.Atoi(idParam)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid character ID")
		return
	}

	comments, err := app.models.Comments.GetCommentsByCharacterID(characterID)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to get character comments")
		return
	}

	app.respondWithJSON(w, http.StatusOK, comments)
}

// выводим список комментов от определенного юзера
func (app *application) getUserCommentsHandler(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	comments, err := app.models.Comments.GetCommentsByUser(userID)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to get comments")
		return
	}

	app.respondWithJSON(w, http.StatusOK, comments)
}
