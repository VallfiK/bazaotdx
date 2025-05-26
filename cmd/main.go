package main

import (
	"guest-management/internal/db"
	"guest-management/internal/service"
	"guest-management/internal/ui"
	"log"
)

func main() {
	// Инициализация БД
	database, err := db.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Конфигурация
	documentsRoot := "documents" // Можно заменить на os.Getenv("DOCUMENTS_PATH")

	// Инициализация сервисов
	guestService := service.NewGuestService(database.DB, documentsRoot)
	cottageService := service.NewCottageService(database.DB)

	// Создание GUI
	app := ui.NewGuestApp(guestService, cottageService)
	app.Run()
}
