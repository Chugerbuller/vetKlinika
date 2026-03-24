package main

import (
	"log"
	"os"
	"vetklinika/backend/internal/server"
	"vetklinika/backend/internal/service"
	"vetklinika/backend/internal/storage"
)

// TODO: добавление питомца при создании клиента, удаление питомца при удалении клиента, добавление страницы для создания питомца, добавление страницы для просмотра питомца, добавление страницы для редактирования питомца и сделать работу с файлами
func main() {
	os.Setenv("CONFIG_PATH", "config.yaml")
	cfg, err := service.MustLoad()
	if err != nil {
		log.Fatal(err.Error())
	}
	storage, err := storage.New(cfg.DbPath)
	if err != nil {
		log.Fatalf("Can`t open db %v", err)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			log.Fatalf("Can`t close db %v", err)
		}
	}()
	server := server.NewServer(storage)
	log.Printf("Server been started on port:%s", cfg.Port)
	if err := server.Start(cfg.Port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
