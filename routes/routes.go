package routes

import (
	project "mrkrab-be/controllers/project"
	transaction "mrkrab-be/controllers/transaction"
	user "mrkrab-be/controllers/user"

	"github.com/gorilla/mux"
)

func Handlers() *mux.Router {
	r := mux.NewRouter()
	// project
	r.HandleFunc("/project", project.NewProject).Methods("POST")
	r.HandleFunc("/project", project.GetProject).Methods("GET")
	r.HandleFunc("/project/{id}", project.GetProject).Methods("GET")
	r.HandleFunc("/project/{id}", project.RemoveProject).Methods("DELETE")
	r.HandleFunc("/project/{id}", project.UpdateProject).Methods("PATCH")

	// transaction
	r.HandleFunc("/transaction", transaction.NewTransaction).Methods("POST")
	r.HandleFunc("/transaction", transaction.GetTransaction).Methods("GET")
	r.HandleFunc("/transaction/{id}", transaction.GetTransaction).Methods("GET")
	r.HandleFunc("/transaction/{id}", transaction.RemoveTransaction).Methods("DELETE")
	r.HandleFunc("/transaction/{id}", transaction.UpdateTransaction).Methods("PATCH")

	// user
	r.HandleFunc("/user", user.NewUser).Methods("POST")
	r.HandleFunc("/user", user.GetUser).Methods("GET")
	r.HandleFunc("/user/{id}", user.GetUser).Methods("GET")
	r.HandleFunc("/user/{id}", user.RemoveUser).Methods("DELETE")
	r.HandleFunc("/user/{id}", user.UpdateUser).Methods("PATCH")

	return r
}
