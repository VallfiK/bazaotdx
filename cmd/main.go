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
		// Автоматическое удаление старых бронирований
		_, err := db.Exec(`
			DELETE FROM lesbaza.bookings 
			WHERE status = $1 AND check_out_date <= NOW() - INTERVAL '24 hours'
		`, models.BookingStatusCancelled)
		if err != nil {
			log.Printf("Error auto-delete old bookings: %v", err)
		}

		// Автоматическое выселение гостей
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

		// Автоматическое удаление старых заселенных бронирований
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

			// Отменяем старые заселенные бронирования
			for _, id := range bookingIDs {
				_, err := db.Exec(`
					UPDATE lesbaza.bookings 
					SET status = $1 
					WHERE booking_id = $2`,
					models.BookingStatusCancelled, id,
				)
				if err == nil {
					// Освобождаем домик
					_, err = db.Exec(`
						UPDATE lesbaza.cottages c
						SET status = 'free'
						FROM lesbaza.bookings b
						WHERE b.booking_id = $1
						AND c.cottage_id = b.cottage_id`,
						id,
					)
					log.Printf("Auto cancelled booking %d", id)
				}
			}
		}

		time.Sleep(1 * time.Hour)
	}
}
