package http

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
	"users/db"
	"users/logger"
	"users/tables"
)

var admin *tables.Admin
var key = []byte("mykey")

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return key, nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {

			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

func GenJwt() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["name"] = "John Doe"
	claims["exp"] = time.Now().Add(time.Minute * 40).Unix()
	tokenString, e := token.SignedString(key)
	if e != nil {
		logger.Log.Println("Пакет http, функция GenJwt. Ошибка подписи токена ключем", e.Error())
		return "", e
	}

	return tokenString, nil
}
func validRequestData(w http.ResponseWriter, r *http.Request) {
	if r.Header["Name"] == nil || r.Header["Password"] == nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	if r.Header["Name"][0] == "" || r.Header["Password"][0] == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
}

func Authorization(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	validRequestData(w, r)

	dataAdmin := tables.Admin{}

	resp := db.GetSingleData(admin, r)

	e := resp.Scan(&dataAdmin.Name, &dataAdmin.Password)
	if e != nil {
		logger.Log.Println("Пакет http. func: getUser. Ошибка: не смогли прочитать ответ БД", e.Error())
	}
	if dataAdmin.Name == "" {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	if dataAdmin.Password != r.Header["Password"][0] || dataAdmin.Name != r.Header["Name"][0] {
		http.Error(w, http.StatusText(401), 401)
		return
	}
	str, e := GenJwt()
	if e != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(w, str)
}

func creatAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	validRequestData(w, r)
	nRows, e := db.InsertData(admin, r)
	if e != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	} else if nRows == 0 {
		http.Error(w, http.StatusText(400), 400)
		return
	}

}

func updateAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		return
	}
	validRequestData(w, r)
	nRows, e := db.EditData(admin, r)
	if e != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	} else if nRows == 0 {
		http.NotFound(w, r)
		return
	}

}

func deleteAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		return
	}
	validRequestData(w, r)
	nRows, e := db.DeleteData(admin, r)
	if e != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	} else if nRows == 0 {
		http.NotFound(w, r)
		return
	}

}
