package tables

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"users/logger"
)

type User struct {
	Name string
	Age  int
	Id   int
}

func (u *User) GetAll() string {
	return `SELECT * FROM "users"`
}

func (u *User) GetSingle(r *http.Request) string {
	params := mux.Vars(r)
	return `SELECT * FROM "users" WHERE "Id"=` + params["id"]
}

func (u *User) InsertSingle(r *http.Request) (string, error) {
	var user User
	e := json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()

	if e != nil {
		logger.Log.Print("Пакет tables. аФайл user.go. Функция insert. Не удалось прочитать тело запроса", e.Error())
		return "", e
	}

	.3str := `INSERT INTO "users" ("Name", "Id", "Age") VALUES (` + "'" + user.Name + "'" + ", " + strconv.Itoa(user.Id) + ", " + strconv.Itoa(user.Age) + `)`
	return str, e
}

func (u *User) EditSingle(r *http.Request) (string, error) {
	var user User
	e := json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()

	if e != nil {
		logger.Log.Print("Пакет tables. аФайл user.go. Функция EditSingle. Не удалось прочитать тело запроса", e.Error())
		return "", e
	}

	str := `UPDATE "users" SET "Name"=` + "'" + user.Name + "', " + `"Age"=` + "'" + strconv.Itoa(user.Age) + "'" + ` WHERE "Id" =` + strconv.Itoa(user.Id)
	return str, e
}

func (u *User) DeleteSingle(r *http.Request) string {
	params := mux.Vars(r)
	str := `DELETE FROM "users" WHERE "Id" =` + params["id"]
	return str
}
