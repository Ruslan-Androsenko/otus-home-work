package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
)

type DataBaseConf struct {
	DbHost string `toml:"db_host"`
	DbPort int    `toml:"db_port"`
	DbUser string `toml:"db_user"`
	DbPass string `toml:"db_pass"`
	DbName string `toml:"db_name"`
}

// NewDBConnection подключение к БД сервиса.
func (storage StorageConf) NewDBConnection() *sql.DB {
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
		config.DbUser, config.DbPass, config.DbHost, config.DbPort, config.DbName)
}
