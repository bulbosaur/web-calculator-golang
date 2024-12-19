package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

func GettingConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := new(Config)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
