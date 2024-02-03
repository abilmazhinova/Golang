package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"tsis1/pkg"
)

func main() {
	log.Println("starting API server")
	router := mux.NewRouter()

	router.HandleFunc("/health-check", pkg.HealthCheck).Methods("GET")
	router.HandleFunc("/students", pkg.GetStudents).Methods("GET")
	router.HandleFunc("/students/{id}", pkg.GetStudentByIdHandler).Methods("GET")

	log.Println("creating routes")

	http.Handle("/", router)

	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
