package main

import (
	_ "github.com/lib/pq"

	//"log"
	//"time"
	"users/db"
	"users/http"
	"users/settings"
	//"users/tables"
	//"fmt"
)

func init() {
	e := settings.InitSettings()
	if e != nil {
		return
	}
	http.InitHttp()
}
func main() {
	defer db.Close()
	e := db.Connect()
	if e != nil {
		return
	}
	http.RunLAS()
}
