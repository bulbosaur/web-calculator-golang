package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Init считывает переменные окружения
func Init() {
	viper.SetDefault("server.ORC_HOST", "localhost")
	viper.SetDefault("server.ORC_PORT", "8080")

	viper.SetDefault("time.TIME_ADDITION_MS", 100)
	viper.SetDefault("time.TIME_SUBTRACTION_MS", 100)
	viper.SetDefault("time.TIME_MULTIPLICATIONS_MS", 100)
	viper.SetDefault("time.TIME_DIVISIONS_MS", 100)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("config file not found in paths: %v", viper.ConfigFileUsed())
		} else {
			log.Printf("error reading config file: %v", err)
		}
		log.Print("default values are set")
	}

	log.Print("config has been successfully initialized")

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
	})
	viper.WatchConfig()
}
