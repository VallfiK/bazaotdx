package service

import (
	"database/sql"
	"fmt"

	"github.com/VallfIK/bazaotdx/internal/models"
)

type CottageService struct {
	db *sql.DB
}

func NewCottageService(db *sql.DB) *CottageService {
	return &CottageService{db: db}
}

func (s *CottageService) GetFreeCottages() ([]models.Cottage, error) {
	rows, err := s.db.Query("SELECT cottage_id, name FROM lesbaza.cottages WHERE status = 'free'")
	if err != nil {
		return nil, fmt.Errorf("failed to query free cottages: %w", err)
	}
	defer rows.Close()

	var cottages []models.Cottage
	for rows.Next() {
		var c models.Cottage
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, fmt.Errorf("failed to scan cottage: %w", err)
		}
		cottages = append(cottages, c)
	}
	return cottages, nil
}

func (s *CottageService) AddCottage(name string) error {
	// Добавляем явное указание всех колонок
	_, err := s.db.Exec(
		"INSERT INTO lesbaza.cottages (name, status) VALUES ($1, 'free')",
		name,
	)
	if err != nil {
		return fmt.Errorf("ошибка добавления домика: %w", err)
	}
	return nil
}

func (s *CottageService) GetAllCottages() ([]models.Cottage, error) {
	rows, err := s.db.Query("SELECT cottage_id, name, status FROM lesbaza.cottages ORDER BY cottage_id")
	if err != nil {
		return nil, fmt.Errorf("failed to query all cottages: %w", err)
	}
	defer rows.Close()

	var cottages []models.Cottage
	for rows.Next() {
		var c models.Cottage
		if err := rows.Scan(&c.ID, &c.Name, &c.Status); err != nil {
			return nil, fmt.Errorf("failed to scan cottage: %w", err)
		}
		cottages = append(cottages, c)
	}
	return cottages, nil
}

// DeleteCottage удаляет домик по ID
func (s *CottageService) DeleteCottage(cottageID int) error {
	// Проверяем, нет ли активных бронирований для этого домика
	var activeBookings int
	err := s.db.QueryRow(`
		SELECT COUNT(*) 
		FROM lesbaza.bookings 
		WHERE cottage_id = $1 
		AND status IN ('booked', 'checked_in')
	`, cottageID).Scan(&activeBookings)

	if err != nil {
		return fmt.Errorf("ошибка проверки активных бронирований: %w", err)
	}

	if activeBookings > 0 {
		return fmt.Errorf("невозможно удалить домик: есть %d активных бронирований", activeBookings)
	}

	// Проверяем, не заселен ли кто-то в домике
	var guestCount int
	err = s.db.QueryRow(`
		SELECT COUNT(*) 
		FROM lesbaza.guests 
		WHERE cottage_id = $1
	`, cottageID).Scan(&guestCount)

	if err != nil {
		return fmt.Errorf("ошибка проверки гостей: %w", err)
	}

	if guestCount > 0 {
		return fmt.Errorf("невозможно удалить домик: в нем есть заселенные гости")
	}

	// Удаляем домик
	result, err := s.db.Exec("DELETE FROM lesbaza.cottages WHERE cottage_id = $1", cottageID)
	if err != nil {
		return fmt.Errorf("ошибка удаления домика: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества удаленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("домик с ID %d не найден", cottageID)
	}

	return nil
}

// UpdateCottageName обновляет название домика
func (s *CottageService) UpdateCottageName(cottageID int, newName string) error {
	result, err := s.db.Exec(
		"UPDATE lesbaza.cottages SET name = $1 WHERE cottage_id = $2",
		newName, cottageID,
	)
	if err != nil {
		return fmt.Errorf("ошибка обновления названия домика: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("домик с ID %d не найден", cottageID)
	}

	return nil
}
