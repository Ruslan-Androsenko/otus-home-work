package main

import (
	"log"

	"github.com/BurntSushi/toml"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/app"
	memorystorage "github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/storage/sql"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.

type Config struct {
	Storage StorageConf
	Logger  LoggerConf
}

type StorageConf struct {
	Type     string
	DataBase DataBaseConf
}

type LoggerConf struct {
	Level string
}

func NewConfig() Config {
	var config Config

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatalf("Can not read config file, err: %v \n", err)
	}

	return config
}

// GetStorage Получить объект хранилища.
func (config Config) GetStorage() app.Storage {
	var storage app.Storage

	switch config.Storage.Type {
	case StorageTypeMemory:
		storage = memorystorage.New()
	case StorageTypeDataBase:
		dbConn = config.Storage.NewDBConnection()
		storage = sqlstorage.New(dbConn, logg)
	default:
		logg.Fatal("Storage type not found")
	}

	return storage
}
