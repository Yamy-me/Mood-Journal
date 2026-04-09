package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	moodJournal "Yamy-Gin"
	"Yamy-Gin/routes"
)

func main() {
	db, errDB := sql.Open("sqlite3", `C:\Users\Atai\Desktop\newGoLang\database.db`)
	if errDB != nil {
		log.Fatalf("[ERROR] can't open db: %v", errDB)
	}
	log.Print("DATABASE CONNECTION: ", db.Ping() == nil)

	server, err := moodJournal.NewServer("8080", routes.New().InitRoutes(db))
	if err != nil {
		log.Fatalf("error in main.go: %v", err)
	}

	// Обязательно логируем ошибку падения сервера!
	if err := server.Run(); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
