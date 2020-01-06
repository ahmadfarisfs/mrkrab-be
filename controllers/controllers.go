package controllers

import (
	"mrkrab-be/controllers/project"
	"net/http"

	"github.com/gorilla/mux"
)

// New .
func New() http.Handler {
	r := mux.NewRouter()
	// project
	r.HandleFunc("/project", project.New).Methods("POST")
	r.HandleFunc("/project", project.Get).Methods("GET")
	r.HandleFunc("/project/{id}", project.Get).Methods("GET")

	r.HandleFunc("/project/{id}", project.Remove).Methods("DELETE")
	r.HandleFunc("/project/{id}", project.Update).Methods("PATCH")
	// transaction

	// user
	r.HandleFunc("/user", user.New).Methods("POST")
	r.HandleFunc("/user/{id}", user.Get).Methods("GET")
	r.HandleFunc("/user/{id}", user.Remove).Methods("DELETE")
	r.HandleFunc("/user/{id}", user.Update).Methods("UPDATE")

	return r
}
