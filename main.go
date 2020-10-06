package main

import (
	"fmt"
	"users/db"
	"users/http"
	"users/settings"
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
	str, e := http.GenJwt()
	fmt.Println(str)
	e = db.Connect()
	if e != nil {
		return
	}
	http.RunLAS()
}
