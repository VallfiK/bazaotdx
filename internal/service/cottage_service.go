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

// internal/service/cottage_service.go
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
	rows, err := s.db.Query("SELECT cottage_id, name, status FROM lesbaza.cottages")
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
