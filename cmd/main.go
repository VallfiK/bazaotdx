package main

import (
	"log"
	"time"

	"github.com/VallfIK/bazaotdx/internal/db"
	"github.com/VallfIK/bazaotdx/internal/service"
	"github.com/VallfIK/bazaotdx/internal/ui"
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
	tariffService := service.NewTariffService(database.DB) // Добавлено

	// Создание GUI
	app := ui.NewGuestApp(guestService, cottageService, tariffService)
	app.Run()

	// Добавьте фоновую задачу
	go func() {
		for {
			_, err := database.DB.Exec(`
            DELETE FROM lesbaza.guests 
            WHERE check_out_date <= NOW() - INTERVAL '2 hours'
        `)
			if err != nil {
				log.Println("Auto-checkout error:", err)
			}
			time.Sleep(1 * time.Hour)
		}
	}()
}
