package db

import (
	"database/sql"
	"fmt"
	"net/http"
	"users/logger"
	"users/settings"
	"users/tables"
)

var db *sql.DB
var stringForRequest string

type actions interface {
	GetAll() string
	GetSingle(*http.Request) string
	InsertSingle(*http.Request) (string, error)
	EditSingle(*http.Request) (string, error)
	DeleteSingle(*http.Request) string
}

func Connect() error {
	var e error
	db, e = sql.Open(settings.Cfg.PgUser, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", settings.Cfg.PgHost, settings.Cfg.PgPort, settings.Cfg.PgUser, settings.Cfg.PgPass, settings.Cfg.PgDB))
	if e != nil {
		logger.Log.Println("Пакет maine. Файл db.go:13 Открыть базы данных не удалось", e.Error())
		return e
	}
	if e = db.Ping(); e != nil {
		logger.Log.Println("Пакет maine. Файл db.go:13 Соединение с базой данных не установлено", e.Error())
		return e
	}
	return nil
}

func GetAllData(a actions) ([]tables.User, error) {
	user := tables.User{}
	var sliceOfRows []tables.User
	stringForRequest = a.GetAll()

	rows, e := db.Query(stringForRequest)
	if e != nil {
		logger.Log.Println("Пакет db. func: GetAllData. Ошибка чтения строк из базы==", e.Error())
		return sliceOfRows, e
	}

	for rows.Next() {
		e = rows.Scan(&user.Name, &user.Id, &user.Age)
		if e != nil {
			logger.Log.Println("Пакет db. func: GetAllData. Ошибка сохранения данных из строки базы в объект user", e.Error())
			return sliceOfRows, e
		}
		sliceOfRows = append(sliceOfRows, user)
	}

	return sliceOfRows, e
}

func GetSingleData(a actions, r *http.Request) (tables.User, error) {
	user := tables.User{}
	stringForRequest = a.GetSingle(r)
	respons := db.QueryRow(stringForRequest)

	e := respons.Scan(&user.Name, &user.Id, &user.Age)
	if e != nil && e != sql.ErrNoRows {
		logger.Log.Println("Пакет db. func: GetSingleData. Ошибка при поиске строки в БД", e.Error())
	}
	return user, e
}

func InsertData(a actions, r *http.Request) (int64, error) {
	stringForRequest, e := a.InsertSingle(r)
	if e != nil {
		return 0, e
	}
	res, e := db.Exec(stringForRequest)
	if e != nil {
		logger.Log.Println("Пакет db. func: InsertData. Ошибка создания записи в БД", e.Error())
		return 0, e
	}
	nrows, e := res.RowsAffected()
	return nrows, e
}

func EditData(a actions, r *http.Request) (int64, error) {
	stringForRequest, e := a.EditSingle(r)
	if e != nil {
		return 0, e
	}

	res, e := db.Exec(stringForRequest)
	if e != nil {
		logger.Log.Println("Пакет db. func: EditData. Ошибка при редактировании строки в БД", e.Error())
		return 0, e
	}
	nrows, e := res.RowsAffected()

	return nrows, e
}

func DeleteData(a actions, r *http.Request) (int64, error) {
	stringForRequest := a.DeleteSingle(r)
	res, e := db.Exec(stringForRequest)
	if e != nil {
		logger.Log.Println("Пакет db. func: DeleteData. Ошибка при удалении строки в БД", e.Error())
		return 0, e
	}
	nrows, e := res.RowsAffected()

	return nrows, e
}

func Close() {
	db.Close()
}
