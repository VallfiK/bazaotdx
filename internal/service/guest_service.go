// internal/service/guest_service.go
package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/VallfIK/bazaotdx/internal/models"
	"github.com/VallfIK/bazaotdx/internal/utils"
)

type GuestService struct {
	db            *sql.DB
	documentsPath string // Путь для хранения документов из конфига
}

func NewGuestService(db *sql.DB, documentsPath string) *GuestService {
	return &GuestService{
		db:            db,
		documentsPath: documentsPath,
	}
}

func (s *GuestService) RegisterGuest(guest models.Guest, cottageID int) error {
	// Создаем папку для документов гостя
	folderPath, err := utils.CreateGuestFolder(s.documentsPath, guest.FullName)
	if err != nil {
		return fmt.Errorf("ошибка создания папки: %w", err)
	}

	// Копируем документ в целевую папку
	if guest.DocumentScanPath != "" {
		newDocPath, err := utils.CopyDocument(guest.DocumentScanPath, folderPath)
		if err != nil {
			return fmt.Errorf("ошибка копирования документа: %w", err)
		}
		guest.DocumentScanPath = newDocPath
	}

	// Начинаем транзакцию
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}

	// Вставляем запись о госте
	_, err = tx.Exec(`
		INSERT INTO lesbaza.guests 
			(full_name, email, phone, cottage_id, document_scan_path, check_in_date, check_out_date, tariff_id) 
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8)`,
		guest.FullName,
		guest.Email,
		guest.Phone,
		cottageID,
		guest.DocumentScanPath,
		guest.CheckInDate,
		guest.CheckOutDate,
		guest.TariffID,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка добавления гостя: %w", err)
	}

	// Обновляем статус домика
	_, err = tx.Exec(`
        UPDATE lesbaza.cottages 
        SET status = 'occupied' 
        WHERE cottage_id = $1`,
		cottageID,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка обновления статуса домика: %w", err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ошибка фиксации транзакции: %w", err)
	}

	return nil
}

// internal/service/guest_service.go
func (s *GuestService) GetGuestByCottageID(cottageID int) (*models.Guest, error) {
	row := s.db.QueryRow(`
        SELECT 
            guest_id, 
            full_name, 
            email, 
            phone, 
            cottage_id, 
            document_scan_path 
        FROM lesbaza.guests 
        WHERE cottage_id = $1`, cottageID)

	guest := &models.Guest{}
	err := row.Scan(
		&guest.ID,
		&guest.FullName,
		&guest.Email,
		&guest.Phone,
		&guest.CottageID,
		&guest.DocumentScanPath,
	)

	if err != nil {
		return nil, fmt.Errorf("guest not found: %w", err)
	}
	return guest, nil
}

func (s *GuestService) CheckOutGuest(cottageID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Удаляем гостя
	_, err = tx.Exec("DELETE FROM lesbaza.guests WHERE cottage_id = $1", cottageID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Обновляем статус домика
	_, err = tx.Exec("UPDATE lesbaza.cottages SET status='free' WHERE cottage_id=$1", cottageID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *GuestService) CalculateCost(checkIn, checkOut time.Time, tariffID int) (float64, error) {
	var price float64
	err := s.db.QueryRow(
		"SELECT price_per_day FROM lesbaza.tariffs WHERE tariff_id = $1",
		tariffID,
	).Scan(&price)
	if err != nil {
		return 0, err
	}

	// Нормализируем даты: заезд в 14:00, выезд в 12:00
	checkInNormalized := time.Date(checkIn.Year(), checkIn.Month(), checkIn.Day(), 14, 0, 0, 0, time.Local)
	checkOutNormalized := time.Date(checkOut.Year(), checkOut.Month(), checkOut.Day(), 12, 0, 0, 0, time.Local)

	// Если выезд в тот же день что и заезд - это минимум 1 день
	if checkInNormalized.Format("2006-01-02") == checkOutNormalized.Format("2006-01-02") {
		return price, nil
	}

	// Рассчитываем количество полных дней
	days := int(checkOutNormalized.Sub(checkInNormalized).Hours() / 24)

	// Если получается 0 или меньше дней, возвращаем стоимость за 1 день
	if days <= 0 {
		return price, nil
	}

	return float64(days) * price, nil
}
