package main

import (
	moodJournal "Yamy-Gin"
	"Yamy-Gin/routes"
	"log"
)

func main() {
	server, err := moodJournal.NewServer("8080", routes.New().InitRoutes())

	if err != nil {
		log.Fatalf("error in main.go: %v", err)
	}

	server.Run()

}