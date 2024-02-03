package pkg

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tsis1/api"
	"github.com/gorilla/mux"
)

func GetStudents(w http.ResponseWriter, r *http.Request) {
	
	Students := api.Students

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResponse, err := json.Marshal(Students)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}

func GetStudentByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]


	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID. Please try again.", http.StatusBadRequest)
		return
	}


	student := api.GetStudentById(id)
	w.Header().Set("Content-Type", "application/json")


	if student == nil {
		http.Error(w, "Student not found. Please try again.", http.StatusNotFound)
		return
	}


	jsonResponse, err := json.Marshal(student)
	if err != nil {
		http.Error(w, "Error encoding JSON.", http.StatusInternalServerError)
		return
	}


	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
