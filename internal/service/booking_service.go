package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/VallfIK/bazaotdx/internal/models"
)

type BookingService struct {
	db            *sql.DB
	tariffService *TariffService
}

func NewBookingService(db *sql.DB) *BookingService {
	return &BookingService{
		db:            db,
		tariffService: NewTariffService(db),
	}
}

// CreateBooking создает новую бронь
func (s *BookingService) CreateBooking(booking models.Booking) (*models.Booking, error) {
	// Проверяем доступность домика на эти даты
	available, err := s.IsCottageAvailable(booking.CottageID, booking.CheckInDate, booking.CheckOutDate)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, fmt.Errorf("домик недоступен на выбранные даты")
	}

	// Получаем тариф
	tariff, err := s.tariffService.GetTariffByID(booking.TariffID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения тарифа: %w", err)
	}

	// Рассчитываем стоимость
	checkInDateOnly := time.Date(booking.CheckInDate.Year(), booking.CheckInDate.Month(), booking.CheckInDate.Day(), 0, 0, 0, 0, time.Local)
	checkOutDateOnly := time.Date(booking.CheckOutDate.Year(), booking.CheckOutDate.Month(), booking.CheckOutDate.Day(), 0, 0, 0, 0, time.Local)
	days := int(checkOutDateOnly.Sub(checkInDateOnly).Hours()/24) + 1
	totalCost := float64(days) * tariff.PricePerDay

	// Создаем бронь
	var bookingID int
	err = s.db.QueryRow(`
		INSERT INTO lesbaza.bookings 
		(cottage_id, guest_name, phone, email, check_in_date, check_out_date, 
		 status, created_at, notes, tariff_id, total_cost)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING booking_id`,
		booking.CottageID,
		booking.GuestName,
		booking.Phone,
		booking.Email,
		booking.CheckInDate,
		booking.CheckOutDate,
		models.BookingStatusBooked,
		time.Now(),
		booking.Notes,
		booking.TariffID,
		totalCost,
	).Scan(&bookingID)

	if err != nil {
		return nil, fmt.Errorf("ошибка создания брони: %w", err)
	}

	booking.ID = bookingID
	booking.Status = models.BookingStatusBooked
	booking.CreatedAt = time.Now()

	return &booking, nil
}

// IsCottageAvailable проверяет доступность домика на даты
func (s *BookingService) IsCottageAvailable(cottageID int, checkIn, checkOut time.Time) (bool, error) {
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM lesbaza.bookings
		WHERE cottage_id = $1
		AND status IN ($2, $3, $4)
		AND NOT (check_out_date <= $5 OR check_in_date >= $6)`,
		cottageID,
		models.BookingStatusBooked,
		models.BookingStatusCheckedIn,
		models.BookingStatusTemporary,
		checkIn,
		checkOut,
	).Scan(&count)

	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// GetBookingsByDateRange получает все брони за период
func (s *BookingService) GetBookingsByDateRange(startDate, endDate time.Time) ([]models.Booking, error) {
	rows, err := s.db.Query(`
		SELECT booking_id, cottage_id, guest_name, phone, email, 
		       check_in_date, check_out_date, status, created_at, notes,
		       COALESCE(tariff_id, 0), COALESCE(total_cost, 0)
		FROM lesbaza.bookings
		WHERE (check_in_date <= $2 AND check_out_date >= $1)
		AND status NOT IN ($3, $4)
		ORDER BY check_in_date`,
		startDate, endDate, models.BookingStatusCancelled, models.BookingStatusCompleted,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var b models.Booking
		err := rows.Scan(
			&b.ID, &b.CottageID, &b.GuestName, &b.Phone, &b.Email,
			&b.CheckInDate, &b.CheckOutDate, &b.Status, &b.CreatedAt, &b.Notes,
			&b.TariffID, &b.TotalCost,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}

	return bookings, nil
}

// GetCalendarData получает данные для календаря
func (s *BookingService) GetCalendarData(startDate, endDate time.Time) (map[time.Time]map[int]models.BookingStatus, error) {
	bookings, err := s.GetBookingsByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Создаем карту дат
	calendar := make(map[time.Time]map[int]models.BookingStatus)

	// Инициализируем все дни
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		calendar[d] = make(map[int]models.BookingStatus)
	}

	// Заполняем данные бронирований
	for _, booking := range bookings {
		// Получаем даты заезда и выезда в начале дня
		checkIn := time.Date(booking.CheckInDate.Year(), booking.CheckInDate.Month(), booking.CheckInDate.Day(), 0, 0, 0, 0, time.Local)
		checkOut := time.Date(booking.CheckOutDate.Year(), booking.CheckOutDate.Month(), booking.CheckOutDate.Day(), 0, 0, 0, 0, time.Local)

		// Проверяем, что даты корректные
		if checkOut.Before(checkIn) {
			continue // Пропускаем некорректные брони
		}

		// Проходим по всем дням брони
		for d := checkIn; d.Before(checkOut); d = d.AddDate(0, 0, 1) {
			if dayMap, exists := calendar[d]; exists {
				isCheckIn := d.Equal(checkIn)
				isCheckOut := d.Equal(checkOut.AddDate(0, 0, -1)) // Последний день проживания

				// Добавляем статус для дня брони
				dayMap[booking.CottageID] = models.BookingStatus{
					Status:     booking.Status,
					BookingID:  booking.ID,
					GuestName:  booking.GuestName,
					IsPartDay:  isCheckIn || isCheckOut,
					IsCheckIn:  isCheckIn,
					IsCheckOut: isCheckOut,
				}
			}
		}

		// Отдельно обрабатываем день выезда (если он в пределах календаря)
		if dayMap, exists := calendar[checkOut]; exists {
			// В день выезда помечаем специальным статусом для диагональной кнопки
			dayMap[booking.CottageID] = models.BookingStatus{
				Status:     booking.Status,
				BookingID:  booking.ID,
				GuestName:  booking.GuestName,
				IsPartDay:  true,
				IsCheckIn:  false,
				IsCheckOut: true,
			}
		}
	}

	return calendar, nil
}

// CheckOutBooking выселяет гостя (завершает бронирование)
func (s *BookingService) CheckOutBooking(bookingID int) error {
	// Получаем бронь
	booking, err := s.GetBookingByID(bookingID)
	if err != nil {
		return err
	}

	if booking.Status != models.BookingStatusCheckedIn {
		return fmt.Errorf("можно выселить только заселенного гостя")
	}

	// Начинаем транзакцию
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Обновляем статус брони на "completed"
	_, err = tx.Exec(
		"UPDATE lesbaza.bookings SET status = $1 WHERE booking_id = $2",
		models.BookingStatusCompleted, bookingID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Освобождаем домик
	_, err = tx.Exec(
		"UPDATE lesbaza.cottages SET status = 'free' WHERE cottage_id = $1",
		booking.CottageID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Удаляем из таблицы гостей, если есть
	_, err = tx.Exec(
		"DELETE FROM lesbaza.guests WHERE cottage_id = $1",
		booking.CottageID,
	)
	if err != nil {
		// Это не критичная ошибка, продолжаем
		// tx.Rollback()
		// return err
	}

	return tx.Commit()
}

// UpdateCheckOutDate обновляет дату выезда (для раннего выселения)
func (s *BookingService) UpdateCheckOutDate(bookingID int, newCheckOutDate time.Time, reason string) error {
	// Получаем бронь
	booking, err := s.GetBookingByID(bookingID)
	if err != nil {
		return err
	}

	if booking.Status != models.BookingStatusCheckedIn {
		return fmt.Errorf("можно изменить дату выезда только для заселенного гостя")
	}

	// Проверяем, что новая дата не раньше даты заезда
	checkInDate := time.Date(booking.CheckInDate.Year(), booking.CheckInDate.Month(), booking.CheckInDate.Day(), 0, 0, 0, 0, time.Local)
	newCheckOutDateOnly := time.Date(newCheckOutDate.Year(), newCheckOutDate.Month(), newCheckOutDate.Day(), 0, 0, 0, 0, time.Local)

	if newCheckOutDateOnly.Before(checkInDate) {
		return fmt.Errorf("дата выезда не может быть раньше даты заезда")
	}

	// Получаем тариф для пересчета стоимости
	tariff, err := s.tariffService.GetTariffByID(booking.TariffID)
	if err != nil {
		return fmt.Errorf("ошибка получения тарифа: %w", err)
	}

	// Пересчитываем стоимость
	days := int(newCheckOutDateOnly.Sub(checkInDate).Hours()/24) + 1
	if days <= 0 {
		days = 1
	}
	newTotalCost := float64(days) * tariff.PricePerDay

	// Формируем примечание
	note := fmt.Sprintf("Изменена дата выезда с %s на %s",
		booking.CheckOutDate.Format("02.01.2006"),
		newCheckOutDate.Format("02.01.2006"))
	if reason != "" {
		note += fmt.Sprintf(". Причина: %s", reason)
	}

	// Обновляем в базе данных
	_, err = s.db.Exec(`
		UPDATE lesbaza.bookings 
		SET check_out_date = $1, total_cost = $2, notes = COALESCE(notes, '') || $3
		WHERE booking_id = $4`,
		newCheckOutDate, newTotalCost, ". "+note, bookingID,
	)
	if err != nil {
		return fmt.Errorf("ошибка обновления даты выезда: %w", err)
	}

	return nil
}

// UpdateBookingStatus обновляет статус брони
func (s *BookingService) UpdateBookingStatus(bookingID int, status string) error {
	_, err := s.db.Exec(
		"UPDATE lesbaza.bookings SET status = $1 WHERE booking_id = $2",
		status, bookingID,
	)
	return err
}

// GetBookingByID получает бронь по ID
func (s *BookingService) GetBookingByID(bookingID int) (*models.Booking, error) {
	var booking models.Booking
	err := s.db.QueryRow(`
		SELECT b.booking_id, b.cottage_id, b.guest_name, b.phone, b.email, 
		       b.check_in_date, b.check_out_date, b.status, b.created_at, 
		       b.notes, b.tariff_id, b.total_cost
		FROM lesbaza.bookings b
		WHERE b.booking_id = $1`,
		bookingID,
	).Scan(
		&booking.ID, &booking.CottageID, &booking.GuestName, &booking.Phone, &booking.Email,
		&booking.CheckInDate, &booking.CheckOutDate, &booking.Status, &booking.CreatedAt,
		&booking.Notes, &booking.TariffID, &booking.TotalCost,
	)
	if err != nil {
		return nil, err
	}

	return &booking, nil
}

// CancelBooking отменяет бронь
func (s *BookingService) CancelBooking(bookingID int) error {
	return s.UpdateBookingStatus(bookingID, models.BookingStatusCancelled)
}

// CheckInBooking заселяет гостя
func (s *BookingService) CheckInBooking(bookingID int) error {
	// Получаем бронь
	booking, err := s.GetBookingByID(bookingID)
	if err != nil {
		return err
	}

	if booking.Status != models.BookingStatusBooked {
		return fmt.Errorf("можно заселить только забронированного гостя")
	}

	// Начинаем транзакцию
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Обновляем статус брони на "заселен"
	_, err = tx.Exec(
		"UPDATE lesbaza.bookings SET status = $1 WHERE booking_id = $2",
		models.BookingStatusCheckedIn, bookingID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Занятым делаем домик
	_, err = tx.Exec(
		"UPDATE lesbaza.cottages SET status = 'occupied' WHERE cottage_id = $1",
		booking.CottageID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Добавляем гостя в таблицу гостей
	_, err = tx.Exec(
		"INSERT INTO lesbaza.guests (cottage_id, full_name, phone, email, check_in_date) VALUES ($1, $2, $3, $4, $5)",
		booking.CottageID, booking.GuestName, booking.Phone, booking.Email, booking.CheckInDate,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// GetAvailableCottagesForDates получает доступные домики на даты
func (s *BookingService) GetAvailableCottagesForDates(checkIn, checkOut time.Time) ([]models.Cottage, error) {
	var cottages []models.Cottage
	rows, err := s.db.Query(`
		SELECT c.cottage_id, c.name, c.status
		FROM lesbaza.cottages c
		WHERE c.cottage_id NOT IN (
			SELECT b.cottage_id
			FROM lesbaza.bookings b
			WHERE b.status NOT IN ('cancelled', 'completed')
			AND (
				($1 BETWEEN b.check_in_date AND b.check_out_date)
				OR ($2 BETWEEN b.check_in_date AND b.check_out_date)
				OR (b.check_in_date BETWEEN $1 AND $2)
			)
		)
	`, checkIn, checkOut)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cottage models.Cottage
		err := rows.Scan(&cottage.ID, &cottage.Name, &cottage.Status)
		if err != nil {
			return nil, err
		}
		cottages = append(cottages, cottage)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cottages, nil
}

// GetUpcomingBookings получает предстоящие брони
func (s *BookingService) GetUpcomingBookings() ([]models.Booking, error) {
	var bookings []models.Booking
	rows, err := s.db.Query(`
		SELECT b.booking_id, b.cottage_id, b.guest_name, b.phone, b.email, 
		       b.check_in_date, b.check_out_date, b.status, b.created_at, 
		       b.notes, b.tariff_id, b.total_cost
		FROM lesbaza.bookings b
		WHERE b.status = 'booked' AND b.check_in_date >= CURRENT_DATE
		ORDER BY b.check_in_date ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(
			&booking.ID, &booking.CottageID, &booking.GuestName, &booking.Phone, &booking.Email,
			&booking.CheckInDate, &booking.CheckOutDate, &booking.Status, &booking.CreatedAt,
			&booking.Notes, &booking.TariffID, &booking.TotalCost,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

// GetBookingsByStatus получает брони по статусу
func (s *BookingService) GetBookingsByStatus(status string, afterDate time.Time) ([]models.Booking, error) {
	var bookings []models.Booking
	rows, err := s.db.Query(`
		SELECT b.booking_id, b.cottage_id, b.guest_name, b.phone, b.email, 
		       b.check_in_date, b.check_out_date, b.status, b.created_at, 
		       b.notes, b.tariff_id, b.total_cost
		FROM lesbaza.bookings b
		WHERE b.status = $1 AND b.check_in_date >= $2
		ORDER BY b.check_in_date ASC
	`, status, afterDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(
			&booking.ID, &booking.CottageID, &booking.GuestName, &booking.Phone, &booking.Email,
			&booking.CheckInDate, &booking.CheckOutDate, &booking.Status, &booking.CreatedAt,
			&booking.Notes, &booking.TariffID, &booking.TotalCost,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}
