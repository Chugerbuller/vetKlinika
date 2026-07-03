package service

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Configuration struct {
	DbPath     string `yaml:"db_path"`
	Port       string `yaml:"port"`
	LoggerPath string `yaml:"logger_path"`
}

func MustLoad() (*Configuration, error) {
	// Получаем путь до конфиг-файла из env-переменной CONFIG_PATH
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("CONFIG_PATH environment variable is not set")
	}

	// Проверяем существование конфиг-файла
	if _, err := os.Stat(configPath); err != nil {
		return nil, fmt.Errorf("error opening config file: %s", err)
	}

	var cfg Configuration

	// Читаем конфиг-файл и заполняем нашу структуру
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %s", err)
	}

	return &cfg, nil
}
