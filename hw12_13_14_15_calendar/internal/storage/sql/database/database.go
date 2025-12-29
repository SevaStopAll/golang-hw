package database

import (
	"time"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

// ConnectDbWithCfg подключиться к базе данных с переданным конфигом.
func ConnectDbWithCfg(dbDriverName string, dsn string) *sqlx.DB {
	DB = sqlx.MustConnect(dbDriverName, dsn)
	// Настройки ниже конфигурируют пулл подключений к базе данных. Их названия стандартны для большинства библиотек.
	// Ознакомиться с их описанием можно на примере документации Hikari pool:
	// https://github.com/brettwooldridge/HikariCP?tab=readme-ov-file#gear-configuration-knobs-baby
	DB.SetMaxIdleConns(5)
	DB.SetMaxOpenConns(20)
	DB.SetConnMaxLifetime(1 * time.Minute)
	DB.SetConnMaxIdleTime(10 * time.Minute)
	return DB
}
