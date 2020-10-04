package http

import (
	"github.com/gorilla/mux"
	"net/http"
	"users/settings"
)

var router = mux.NewRouter()

//var srv
func InitHttp() {
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users", creatUser).Methods("POST")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
}

func RunLAS() {
	http.ListenAndServe(settings.Cfg.ServerHost+":"+settings.Cfg.ServerPort, router)
}
