package config

import (
	"database/sql"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/queue/rabbitmq"
	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/server"
	memorystorage "github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/storage/sql"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.

type Config struct {
	Server  server.Conf
	Storage StorageConf
	Logger  LoggerConf
	Queue   QueueConf
}

type StorageConf struct {
	Type     string
	DataBase DataBaseConf
}

type LoggerConf struct {
	Level string
}

var dbConn *sql.DB

func NewConfig(configFile string) Config {
	var config Config

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatalf("Can not read config file, err: %v \n", err)
	}

	return config
}

// GetStorage Получить объект хранилища.
func (config Config) GetStorage(logg *logger.Logger) app.Storage {
	var storage app.Storage

	switch config.Storage.Type {
	case StorageTypeMemory:
		storage = memorystorage.New()
	case StorageTypeDataBase:
		dbConn = config.Storage.NewDBConnection(logg)
		storage = sqlstorage.New(dbConn, logg)
	default:
		logg.Fatal("Storage type not found")
	}

	return storage
}

// GetQueue Получить объект очереди.
func (config Config) GetQueue() Queue {
	queueConf := config.Queue
	dataSource := queueConf.getDataSource()

	return rabbitmq.New(dataSource, queueConf.Exchange)
}
