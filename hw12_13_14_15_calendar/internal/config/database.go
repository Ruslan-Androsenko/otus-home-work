//nolint:revive
package config

import (
	"database/sql"
	"fmt"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
)

type DataBaseConf struct {
	DBHost string `toml:"db_host"`
	DBPort int    `toml:"db_port"`
	DBUser string `toml:"db_user"`
	DBPass string `toml:"db_pass"`
	DBName string `toml:"db_name"`
}

// NewDBConnection подключение к БД сервиса.
func (storage StorageConf) NewDBConnection(logg *logger.Logger) *sql.DB {
	dataSource := storage.DataBase.getDataSource()
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		logg.Fatalf("Can't open database connection. Error: %v", err)
	}

	err = goose.SetDialect("mysql")
	if err != nil {
		logg.Fatalf("Failed to set SQL dialect. Error: %v", err)
	}

	return db
}

// Получить строку соединения для БД.
func (config DataBaseConf) getDataSource() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName)
}
