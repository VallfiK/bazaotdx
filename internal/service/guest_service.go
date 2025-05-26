// internal/service/guest_service.go
package service

import (
	"database/sql"
	"fmt"

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
            (full_name, email, phone, cottage_id, document_scan_path) 
        VALUES 
            ($1, $2, $3, $4, $5)`,
		guest.FullName,
		guest.Email,
		guest.Phone,
		cottageID,
		guest.DocumentScanPath,
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
