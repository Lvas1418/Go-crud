package settings

import (
	"encoding/json"
	"os"
	"users/logger"
)

type settings struct {
	ServerHost string
	ServerPort string
	PgHost     string
	PgPort     string
	PgUser     string
	PgPass     string
	PgDB       string
}

var Cfg *settings

func InitSettings() error {
	var e error
	file, e := os.Open("settings/settings.cfg")
	if e != nil {
		logger.Log.Println("Пакет settings. Файл settings.go. Не удалось открыит файл settings.cfg", e.Error())
		return e
	}

	defer file.Close()

	stat, e := file.Stat()
	if e != nil {
		logger.Log.Println("Пакет settings. Файл settings.go. Информация о файле settings.cfg не получена", e.Error())
		return e
	}

	bytes := make([]byte, stat.Size())
	_, e = file.Read(bytes)
	if e != nil {
		logger.Log.Println("Пакет settings. Файл settings.go. Файл settings.cfg не прочитан", e.Error())
		return e
	}

	_ = json.Unmarshal(bytes, &Cfg)
	return e

}
