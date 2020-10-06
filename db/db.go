package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"users/logger"
	"users/settings"
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
		logger.Log.Println("Пакет maine. Файл db.go: Подключиться не удалось", e.Error())
		return e
	}

	if e = db.Ping(); e != nil {
		logger.Log.Println("Пакет maine. Файл db.go Соединение с базой данных не установлено", e.Error())
		return e
	}

	return nil
}

func GetAllData(a actions) (*sql.Rows, error) {
	stringForRequest = a.GetAll()
	rows, e := db.Query(stringForRequest)
	if e != nil {
		logger.Log.Println("Пакет db. func: GetAllData. Ошибка чтения строк из базы==", e.Error())
		return rows, e
	}
	return rows, e
}

func GetSingleData(a actions, r *http.Request) *sql.Row {

	stringForRequest = a.GetSingle(r)
	respons := db.QueryRow(stringForRequest)
	return respons
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
