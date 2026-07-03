package service

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger(path string) *logrus.Logger {
	var logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Info("Не удалось открыть файл логов, используется стандартный stderr")
	} else {
		logger.SetOutput(file)

	}
	logger.Info("Логгер инициализирован")

	return logger
}
