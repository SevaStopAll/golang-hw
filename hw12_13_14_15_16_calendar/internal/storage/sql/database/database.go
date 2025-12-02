package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

// Временная переменная, которая будет ссылаться на подключение к базе данных. Позже мы от неё избавимся
var DB *sqlx.DB

// ConnectDbWithCfg подключиться к базе данных с переданным конфигом
func ConnectDbWithCfg(DbDriverName string, Dsn string) *sqlx.DB {
	DB = sqlx.MustConnect(DbDriverName, Dsn)
	// Настройки ниже конфигурируют пулл подключений к базе данных. Их названия стандартны для большинства библиотек.
	// Ознакомиться с их описанием можно на примере документации Hikari pool:
	// https://github.com/brettwooldridge/HikariCP?tab=readme-ov-file#gear-configuration-knobs-baby
	DB.SetMaxIdleConns(5)
	DB.SetMaxOpenConns(20)
	DB.SetConnMaxLifetime(1 * time.Minute)
	DB.SetConnMaxIdleTime(10 * time.Minute)
	return DB
}
