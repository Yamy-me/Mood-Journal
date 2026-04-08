package main

import (
	"log"

	moodJournal "Yamy-Gin"
	"Yamy-Gin/routes"
)

func main() {
    server, err := moodJournal.NewServer("8080", routes.New().InitRoutes())
    if err != nil {
        log.Fatalf("error in main.go: %v", err)
    }

    // Обязательно логируем ошибку падения сервера!
    if err := server.Run(); err != nil {
        log.Fatalf("Ошибка при запуске сервера: %v", err)
    }
}