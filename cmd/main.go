package main

import (
	"fmt"
	"log"
	"os"

	"net-cat/internal/models"
)

func main() {
	arguments := os.Args
	if len(arguments) >= 3 {
		fmt.Println("[USAGE]: ./cmd/main.go $port")
		return
	}

	var serverPort string
	if len(arguments) == 1 {
		serverPort = ":8989"
	} else {
		serverPort = ":" + arguments[1]
	}

	// Создание и запуск сервера чата
	chatServer := models.CreateChatServer(serverPort)
	go chatServer.DistributeChatMessages() // Запуск горутины для отправки сообщений
	fmt.Println("Listening on the port " + serverPort)

	// Запуск сервера и логирование возможных ошибок
	log.Fatal(chatServer.Launch())
}
