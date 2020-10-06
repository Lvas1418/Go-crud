package http

import (
	"github.com/gorilla/mux"
	"net/http"
	"users/settings"
)

var router = mux.NewRouter()

//var srv
func InitHttp() {
	router.HandleFunc("/", Authorization).Methods("GET")
	router.HandleFunc("/admin", creatAdmin).Methods("POST")
	router.HandleFunc("/admin", isAuthorized(updateAdmin)).Methods("PUT")
	router.HandleFunc("/admin", isAuthorized(deleteAdmin)).Methods("DELETE")
	router.HandleFunc("/users", isAuthorized(getUsers)).Methods("GET")
	router.HandleFunc("/users/{id}", isAuthorized(getUser)).Methods("GET")
	router.HandleFunc("/users", isAuthorized(creatUser)).Methods("POST")
	router.HandleFunc("/users/{id}", isAuthorized(updateUser)).Methods("PUT")
	router.HandleFunc("/users/{id}", isAuthorized(deleteUser)).Methods("DELETE")
}

func RunLAS() {
	http.ListenAndServe(settings.Cfg.ServerHost+":"+settings.Cfg.ServerPort, router)
}
