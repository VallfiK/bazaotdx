// Измените ваш cmd/main.go следующим образом:

package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/VallfIK/bazaotdx/internal/app"
	"github.com/VallfIK/bazaotdx/internal/db"
	"github.com/VallfIK/bazaotdx/internal/models"
	"github.com/VallfIK/bazaotdx/internal/service"
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
	tariffService := service.NewTariffService(database.DB)
	bookingService := service.NewBookingService(database.DB) // Новый сервис

	// Создание приложения
	app := app.NewGuestApp(guestService, cottageService, tariffService, bookingService)

	// Запускаем фоновые задачи ДО запуска GUI
	go backgroundTasks(database.DB, bookingService)

	// Запуск приложения
	app.Run()
}

// backgroundTasks выполняет фоновые задачи
func backgroundTasks(db *sql.DB, bookingService *service.BookingService) {
	for {
		// Автоматическое удаление старых отмененных и завершенных бронирований (старше 30 дней)
		_, err := db.Exec(`
			DELETE FROM lesbaza.bookings 
			WHERE status IN ($1, $2) AND created_at <= NOW() - INTERVAL '30 days'
		`, models.BookingStatusCancelled, models.BookingStatusCompleted)
		if err != nil {
			log.Printf("Error auto-delete old bookings: %v", err)
		}

		// Автоматическое выселение гостей из таблицы guests
		_, err = db.Exec(`
			DELETE FROM lesbaza.guests 
			WHERE check_out_date <= NOW() - INTERVAL '2 hours'
		`)
		if err != nil {
			log.Println("Auto-checkout error:", err)
		}

		// Автоматическое обновление статусов бронирований
		// Переводим "забронировано" в "заселено" если наступил день заезда
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
					log.Printf("Auto check-in error for booking %d: %v", id, err)
				} else {
					log.Printf("Auto checked-in booking %d", id)
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
					log.Printf("Auto checkout error for booking %d: %v", id, err)
				} else {
					log.Printf("Auto checked out booking %d", id)
				}
			}
		}

		time.Sleep(1 * time.Hour)
	}
}
