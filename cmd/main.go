// cmd/main.go - Обновленная версия с новой темой
package main

import (
	"io"
	"log"
	"os"

	"fyne.io/fyne/v2/app"
	"github.com/VallfIK/bazaotdx/internal/app"
	"github.com/VallfIK/bazaotdx/internal/db"
	"github.com/VallfIK/bazaotdx/internal/service"
	"github.com/VallfIK/bazaotdx/internal/ui"
)

func main() {
	// Настройка логирования
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("❌ Ошибка при создании файла логов: %v", err)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	log.Println("🚀 Запуск приложения с новой темой")

	// Инициализация БД
	database, err := db.NewPostgresDB()
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к БД: %v", err)
	}
	defer database.Close()

	// Конфигурация
	documentsRoot := "documents"

	// Инициализация сервисов
	guestService := service.NewGuestService(database.DB, documentsRoot)
	cottageService := service.NewCottageService(database.DB)
	tariffService := service.NewTariffService(database.DB)
	bookingService := service.NewBookingService(database.DB)

	// Создание приложения с кастомной темой
	fyneApp := app.New()
	fyneApp.Settings().SetTheme(&ui.ResortTheme{}) // Применяем нашу тему

	// Создание главного приложения
	guestApp := app.NewGuestApp(guestService, cottageService, tariffService, bookingService)
	guestApp.SetFyneApp(fyneApp) // Передаем настроенное Fyne приложение

	// Запускаем фоновые задачи
	go backgroundTasks(database.DB, bookingService)

	// Запуск приложения
	guestApp.Run()
}
