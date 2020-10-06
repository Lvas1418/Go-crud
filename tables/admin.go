package tables

import (
	"net/http"
	"users/logger"
)

type Admin struct {
	Id       int
	Name     string
	Password string
}

func (u *Admin) GetAll() string {
	return ``
}

func (u *Admin) GetSingle(r *http.Request) string {
	logger.Log.Print("Request==", `SELECT * FROM "admin" WHERE "Name"=`+"'"+r.Header["Name"][0]+"'")
	return `SELECT * FROM "admin" WHERE "Name"=` + "'" + r.Header["Name"][0] + "'"
}

func (u *Admin) InsertSingle(r *http.Request) (string, error) {
	str := `INSERT INTO "admin" ("Name", "Password") VALUES (` + "'" + r.Header["Name"][0] + "'" + ", " + "'" + r.Header["Password"][0] + "'" + `)`
	return str, nil
}

func (u *Admin) EditSingle(r *http.Request) (string, error) {
	str := `UPDATE "admin" SET "Password"=` + "'" + r.Header["Password"][0] + "'" + ` WHERE "Name" =` + "'" + r.Header["Name"][0] + "'"
	return str, nil
}

func (u *Admin) DeleteSingle(r *http.Request) string {
	str := `DELETE FROM "admin" WHERE "Name" =` + "'" + r.Header["Name"][0] + "'"
	return str
}
