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
	dataUser := tables.User{}
	var sliceOfRows []tables.User

	rows, e := db.GetAllData(user)
	if e == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if e != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	for rows.Next() {
		e = rows.Scan(&dataUser.Id, &dataUser.Name, &dataUser.Age)
		if e != nil {
			logger.Log.Println("Пакет db. func: GetAllData. Ошибка сохранения данных из строки базы в объект user", e.Error())
			return
		}
		sliceOfRows = append(sliceOfRows, dataUser)
	}
	json.NewEncoder(w).Encode(sliceOfRows)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	dataUser := tables.User{}

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
	resp /*, e*/ := db.GetSingleData(user, r)
	/*if e == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if e != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}*/

	e = resp.Scan(&dataUser.Id, &dataUser.Name, &dataUser.Age)
	if e != nil {
		logger.Log.Println("Пакет http. func: getUser. Ошибка: не смогли прочитать ответ БД", e.Error())
	}
	//We launch a goroutine that writes a user to the cache,
	// falls asleep for 10 seconds, and when it wakes up, it deletes the user from the cache
	if dataUser.Id == 0 {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	go cash.InsertAndDEl(dataUser)

	logger.Log.Println("The user was taken from the database")
	json.NewEncoder(w).Encode(dataUser)
}

func creatUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	nRows, e := db.InsertData(user, r)
	if e != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	} else if nRows == 0 {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	getUsers(w, r)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		return
	}

	nRows, e := db.EditData(user, r)
	if e != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	} else if nRows == 0 {
		http.NotFound(w, r)
		return
	}

	getUsers(w, r)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		return
	}

	nRows, e := db.DeleteData(user, r)
	if e != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	} else if nRows == 0 {
		http.NotFound(w, r)
		return
	}

	getUsers(w, r)
}
