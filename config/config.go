package config

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Init() считывает переменные окружения
func Init() {
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", "8080")

	viper.SetDefault("time.TIME_ADDITION_MS", 100)
	viper.SetDefault("time.TIME_SUBTRACTION_MS", 100)
	viper.SetDefault("time.TIME_MULTIPLICATIONS_MS", 100)
	viper.SetDefault("time.TIME_DIVISIONS_MS", 100)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Print("config file not found, default values are set")
		} else {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	log.Print("config has been successfully initialized")

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
	})
	viper.WatchConfig()
}
