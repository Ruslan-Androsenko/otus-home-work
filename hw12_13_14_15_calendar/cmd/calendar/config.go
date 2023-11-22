package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf
	// TODO
}

type LoggerConf struct {
	Level string
	// TODO
}

func NewConfig() Config {
	var config Config

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatalf("Can not read config file, err: %v \n", err)
	}

	return config
}

// TODO
