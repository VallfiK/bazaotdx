package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/VallfIK/bazaotdx/internal/models"
)

type BookingService struct {
	db *sql.DB
}

func NewBookingService(db *sql.DB) *BookingService {
	return &BookingService{db: db}
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
		booking.TotalCost,
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
		AND status != $3
		ORDER BY check_in_date`,
		startDate, endDate, models.BookingStatusCancelled,
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
			return nil, fmt.Errorf("некорректные даты бронирования: дата выезда раньше даты заезда")
		}

		// Проходим по всем дням брони
		for d := checkIn; !d.After(checkOut); d = d.AddDate(0, 0, 1) {
			if _, exists := calendar[d]; exists {
				isCheckIn := d.Equal(checkIn)
				isCheckOut := d.Equal(checkOut)

				calendar[d][booking.CottageID] = models.BookingStatus{
					Status:     booking.Status,
					BookingID:  booking.ID,
					GuestName:  booking.GuestName,
					IsPartDay:  isCheckIn || isCheckOut,
					IsCheckIn:  isCheckIn,
					IsCheckOut: isCheckOut,
				}
			}
		}
	}

	return calendar, nil
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
	var b models.Booking
	err := s.db.QueryRow(`
		SELECT booking_id, cottage_id, guest_name, phone, email, 
		       check_in_date, check_out_date, status, created_at, notes,
		       COALESCE(tariff_id, 0), COALESCE(total_cost, 0)
		FROM lesbaza.bookings
		WHERE booking_id = $1`, bookingID,
	).Scan(
		&b.ID, &b.CottageID, &b.GuestName, &b.Phone, &b.Email,
		&b.CheckInDate, &b.CheckOutDate, &b.Status, &b.CreatedAt, &b.Notes,
		&b.TariffID, &b.TotalCost,
	)

	if err != nil {
		return nil, err
	}

	return &b, nil
}

// CancelBooking отменяет бронь
func (s *BookingService) CancelBooking(bookingID int) error {
	return s.UpdateBookingStatus(bookingID, models.BookingStatusCancelled)
}

// CheckOutBooking выселение гостя
func (s *BookingService) CheckOutBooking(bookingID int) error {
	// Получаем текущую бронь
	booking, err := s.GetBookingByID(bookingID)
	if err != nil {
		return fmt.Errorf("ошибка получения брони: %w", err)
	}

	// Проверяем статус брони
	if booking.Status != models.BookingStatusCheckedIn {
		return fmt.Errorf("бронь не имеет статуса " + models.BookingStatusCheckedIn)
	}

	// Обновляем статус на "Выселено"
	return s.UpdateBookingStatus(bookingID, models.BookingStatusCheckedOut)
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

	// Обновляем статус брони
	_, err = tx.Exec(
		"UPDATE lesbaza.bookings SET status = $1 WHERE booking_id = $2",
		models.BookingStatusCheckedIn, bookingID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Обновляем статус домика
	_, err = tx.Exec(
		"UPDATE lesbaza.cottages SET status = 'occupied' WHERE cottage_id = $1",
		booking.CottageID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// GetAvailableCottagesForDates получает доступные домики на даты
func (s *BookingService) GetAvailableCottagesForDates(checkIn, checkOut time.Time) ([]models.Cottage, error) {
	rows, err := s.db.Query(`
		SELECT c.cottage_id, c.name, c.status
		FROM lesbaza.cottages c
		WHERE c.cottage_id NOT IN (
			SELECT DISTINCT cottage_id 
			FROM lesbaza.bookings
			WHERE status IN ($1, $2)
			AND NOT (check_out_date <= $3 OR check_in_date >= $4)
		)
		ORDER BY c.name`,
		models.BookingStatusBooked,
		models.BookingStatusCheckedIn,
		checkIn,
		checkOut,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cottages []models.Cottage
	for rows.Next() {
		var c models.Cottage
		err := rows.Scan(&c.ID, &c.Name, &c.Status)
		if err != nil {
			return nil, err
		}
		cottages = append(cottages, c)
	}

	return cottages, nil
}

// GetUpcomingBookings получает предстоящие брони
func (s *BookingService) GetUpcomingBookings() ([]models.Booking, error) {
	today := time.Now().Truncate(24 * time.Hour)
	return s.GetBookingsByStatus(models.BookingStatusBooked, today)
}

// GetBookingsByStatus получает брони по статусу
func (s *BookingService) GetBookingsByStatus(status string, afterDate time.Time) ([]models.Booking, error) {
	rows, err := s.db.Query(`
		SELECT booking_id, cottage_id, guest_name, phone, email, 
		       check_in_date, check_out_date, status, created_at, notes,
		       COALESCE(tariff_id, 0), COALESCE(total_cost, 0)
		FROM lesbaza.bookings
		WHERE status = $1 AND check_in_date >= $2
		ORDER BY check_in_date`,
		status, afterDate,
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
