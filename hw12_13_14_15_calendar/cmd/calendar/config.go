package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger       LoggerConf
	DbDriverName string
	Dsn          string
	Host         string
	Port         int
	Level        string
	Location     string `toml:"LogLocation"`
	Storage      string
}

type LoggerConf struct {
	Level    string
	Location string
}

func NewConfig(configFile string) Config {
	var config Config
	_, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		fmt.Println(err)
	}
	config.Logger.Level = config.Level
	config.Logger.Location = config.Location

	return config
}
