package http

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"users/cash"
	"users/db"
	"users/logger"
	"users/tables"
)

var user *tables.User

func getUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	resp, e := db.GetAllData(user)
	if e == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if e != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	//first we look in the cache,
	params := mux.Vars(r)
	str := params["id"]
	id, e := strconv.Atoi(str)
	if e != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	result, ok := cash.Find(id)
	if ok {
		logger.Log.Println("The user was taken from the cache")
		json.NewEncoder(w).Encode(result)
		return
	}
	//if the user is not found in the cache, then we go to the database
	resp, e := db.GetSingleData(user, r)
	if e == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if e != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	//We launch a goroutine that writes a user to the cache,
	// falls asleep for 10 seconds, and when it wakes up, it deletes the user from the cache
	go cash.InsertAndDEl(resp)

	logger.Log.Println("The user was taken from the database")
	json.NewEncoder(w).Encode(resp)
}

func creatUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	w.Header().Set("Content-Type", "application/json")

	nRows, e := db.InsertData(user, r)
	if e != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	} else if nRows == 0 {
		http.NotFound(w, r)
		return
	}

	resp, e := db.GetAllData(user)
	if e != nil {
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	nRows, e := db.EditData(user, r)
	if e != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	} else if nRows == 0 {
		http.NotFound(w, r)
		return
	}

	data, e := db.GetAllData(user)
	if e != nil {
		return
	}

	json.NewEncoder(w).Encode(data)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	nRows, e := db.DeleteData(user, r)
	if e != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	} else if nRows == 0 {
		http.NotFound(w, r)
		return
	}

	data, e := db.GetAllData(user)
	if e != nil {
		return
	}

	json.NewEncoder(w).Encode(data)
}
