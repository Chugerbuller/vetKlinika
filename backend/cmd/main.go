package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"vetklinika/backend/internal/server"
	"vetklinika/backend/internal/service"
	"vetklinika/backend/internal/storage"

	"github.com/sirupsen/logrus"
)

// TODO: добавление питомца при создании клиента, удаление питомца при удалении клиента, добавление страницы для создания питомца, добавление страницы для просмотра питомца, добавление страницы для редактирования питомца и сделать работу с файлами
func main() {
	os.Setenv("CONFIG_PATH", "config.yaml")
	cfg, err := service.MustLoad()

	if err != nil {
		log.Fatal(err.Error())
	}
	logger := service.NewLogger(cfg.LoggerPath)
	storage, err := storage.New(cfg.DbPath)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Can`t open db")
	}
	defer func() {
		if err := storage.Close(); err != nil {
			logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatal("Can`t close db")
		}
	}()

	server := server.NewServer(storage, logger)
	url := fmt.Sprintf("http://localhost:%s", cfg.Port)
	ready := make(chan struct{})

	go func() {
		<-ready // Ждем сигнала (закрытия канала)
		log.Printf("Сервер готов! Открываем браузер: %s", url)
		openBrowser(url)
	}()
	log.Printf("Server been started on port:%s", cfg.Port)

	if err := server.Start(cfg.Port, ready); err != nil {
		log.Fatalf("Server error: %v", err)
	}

}
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		// Используем rundll32 для надежности в Windows
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin": // macOS
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Printf("Не удалось открыть браузер: %v", err)
	}
}
