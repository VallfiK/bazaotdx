// cmd/main.go
package main

import (
	"database/sql"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/VallfIK/bazaotdx/internal/app"
	"github.com/VallfIK/bazaotdx/internal/db"
	"github.com/VallfIK/bazaotdx/internal/models"
	"github.com/VallfIK/bazaotdx/internal/service"
)

func main() {
	// Настройка логирования
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("❌ Ошибка при создании файла логов: %v", err)
	}
	defer logFile.Close()

	// Перенаправляем логи в файл и консоль
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	log.Println("🌲 Запуск системы управления базой отдыха 'Звуки Леса'")
	log.Printf("🛠️ Версия: 3.0 (Улучшенный дизайн)")
	log.Printf("🎨 Тема: Природная (зеленые тона, золотые акценты, кремовый фон)")
	log.Printf("🏞️ База отдыха: Звуки Леса")

	// Проверяем существование файлов изображений
	projectPath := "C:\\Users\\VallfIK\\Documents\\GitHub\\bazaotdx"
	imagesPath := filepath.Join(projectPath, "images")
	log.Printf("🔍 Проверка директории images: %s", imagesPath)

	images := []string{
		"free.png",
		"booked.png",
		"bought.png",
		"freefirst.png",
		"bookedfirst.png",
		"boughtfirst.png",
		"bookedlast.png",
		"boughtlast.png",
	}

	for _, img := range images {
		imgPath := filepath.Join(imagesPath, img)
		if _, err := os.Stat(imgPath); err != nil {
			log.Printf("⚠️ Предупреждение: Файл не найден: %s", imgPath)
			// Создаем отсутствующие изображения
			createMissingImage(imgPath)
		} else {
			log.Printf("✅ Файл найден: %s", imgPath)
		}
	}

	// Инициализация БД
	database, err := db.NewPostgresDB()
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к БД: %v", err)
	} else {
		log.Println("✅ Успешное подключение к БД")
	}
	defer database.Close()

	// Конфигурация
	documentsRoot := "documents"

	// Инициализация сервисов
	guestService := service.NewGuestService(database.DB, documentsRoot)
	cottageService := service.NewCottageService(database.DB)
	tariffService := service.NewTariffService(database.DB)
	bookingService := service.NewBookingService(database.DB)

	// Создание улучшенного приложения "Звуки Леса"
	app := app.NewStyledGuestApp(guestService, cottageService, tariffService, bookingService)

	// Запускаем фоновые задачи
	go backgroundTasks(database.DB, bookingService)

	log.Println("🌲 Запуск системы управления 'Звуки Леса'...")
	log.Println("🎯 Особенности новой версии:")
	log.Println("   • Фиксированные размеры окна (предотвращение расширения)")
	log.Println("   • Улучшенная цветовая схема с природными оттенками")
	log.Println("   • Обновленное название: 'Звуки Леса'")
	log.Println("   • Эмоджи и улучшенная типографика")
	log.Println("   • Градиентные эффекты и стилизованные карточки")

	// Запуск приложения
	app.Run()
}

// createMissingImage создает простое изображение-заглушку
func createMissingImage(path string) {
	// Создаем директорию если не существует
	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0755)

	// Создаем пустой файл как заглушку
	file, err := os.Create(path)
	if err != nil {
		log.Printf("❌ Не удалось создать файл-заглушку: %v", err)
		return
	}
	file.Close()
	log.Printf("✅ Создан файл-заглушка: %s", path)
}

// backgroundTasks выполняет фоновые задачи
func backgroundTasks(db *sql.DB, bookingService *service.BookingService) {
	log.Println("🔄 Запуск фоновых задач для 'Звуки Леса'...")

	for {
		// Автоматическое удаление старых отмененных и завершенных бронирований
		result, err := db.Exec(`
			DELETE FROM lesbaza.bookings 
			WHERE status IN ($1, $2) AND created_at <= NOW() - INTERVAL '30 days'
		`, models.BookingStatusCancelled, models.BookingStatusCompleted)
		if err != nil {
			log.Printf("⚠️ Ошибка автоудаления старых броней: %v", err)
		} else {
			if rowsAffected, _ := result.RowsAffected(); rowsAffected > 0 {
				log.Printf("🗑️ Автоматически удалено %d старых бронирований", rowsAffected)
			}
		}

		// Автоматическое выселение гостей
		result, err = db.Exec(`
			DELETE FROM lesbaza.guests 
			WHERE check_out_date <= NOW() - INTERVAL '2 hours'
		`)
		if err != nil {
			log.Printf("⚠️ Ошибка автовыселения: %v", err)
		} else {
			if rowsAffected, _ := result.RowsAffected(); rowsAffected > 0 {
				log.Printf("🚪 Автоматически выселено %d гостей", rowsAffected)
			}
		}

		// Автоматическое обновление статусов бронирований
		rows, err := db.Query(`
			SELECT booking_id 
			FROM lesbaza.bookings 
			WHERE status = $1 
			AND check_in_date::date = CURRENT_DATE
			AND CURRENT_TIME > TIME '14:00'`,
			models.BookingStatusBooked,
		)
		if err == nil {
			var bookingIDs []int
			for rows.Next() {
				var id int
				if rows.Scan(&id) == nil {
					bookingIDs = append(bookingIDs, id)
				}
			}
			rows.Close()

			// Обновляем статусы
			for _, id := range bookingIDs {
				err := bookingService.CheckInBooking(id)
				if err != nil {
					log.Printf("⚠️ Ошибка автозаселения брони %d: %v", id, err)
				} else {
					log.Printf("✅ Автоматически заселена бронь %d", id)
				}
			}
		}

		// Автоматическое завершение просроченных заселенных бронирований
		rows, err = db.Query(`
			SELECT booking_id 
			FROM lesbaza.bookings 
			WHERE status = $1 
			AND check_out_date::date < CURRENT_DATE`,
			models.BookingStatusCheckedIn,
		)
		if err == nil {
			var bookingIDs []int
			for rows.Next() {
				var id int
				if rows.Scan(&id) == nil {
					bookingIDs = append(bookingIDs, id)
				}
			}
			rows.Close()

			// Выселяем просроченные брони
			for _, id := range bookingIDs {
				err := bookingService.CheckOutBooking(id)
				if err != nil {
					log.Printf("⚠️ Ошибка автовыселения брони %d: %v", id, err)
				} else {
					log.Printf("✅ Автоматически выселена просроченная бронь %d", id)
				}
			}
		}

		// Вывод статистики работы системы каждые 6 часов
		currentTime := time.Now()
		if currentTime.Hour()%6 == 0 && currentTime.Minute() < 2 {
			log.Printf("📊 Система 'Звуки Леса' работает стабильно. Время: %s", currentTime.Format("15:04:05 02.01.2006"))
		}

		// Пауза между проверками
		time.Sleep(1 * time.Hour)
	}
}
