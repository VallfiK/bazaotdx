// Исправления для internal/service/tariff_service.go

package service

import (
	"database/sql"
	"fmt"

	"github.com/VallfIK/bazaotdx/internal/models"
)

type TariffService struct {
	db *sql.DB
}

func NewTariffService(db *sql.DB) *TariffService {
	return &TariffService{db: db}
}

func (s *TariffService) CreateTariff(name string, price float64) error {
	_, err := s.db.Exec(
		"INSERT INTO lesbaza.tariffs (name, price_per_day) VALUES ($1, $2)",
		name, price,
	)
	if err != nil {
		return fmt.Errorf("ошибка создания тарифа: %w", err)
	}
	return nil
}

func (s *TariffService) GetTariffs() ([]models.Tariff, error) {
	rows, err := s.db.Query("SELECT tariff_id, name, price_per_day FROM lesbaza.tariffs ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("ошибка получения тарифов: %w", err)
	}
	defer rows.Close()

	var tariffs []models.Tariff
	for rows.Next() {
		var t models.Tariff
		err := rows.Scan(&t.ID, &t.Name, &t.PricePerDay)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования тарифа: %w", err)
		}
		tariffs = append(tariffs, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке строк: %w", err)
	}

	return tariffs, nil
}

func (s *TariffService) GetTariffByID(tariffID int) (*models.Tariff, error) {
	row := s.db.QueryRow("SELECT tariff_id, name, price_per_day FROM lesbaza.tariffs WHERE tariff_id = $1", tariffID)

	var t models.Tariff
	err := row.Scan(&t.ID, &t.Name, &t.PricePerDay)
	if err != nil {
		return nil, fmt.Errorf("тариф с ID %d не найден: %w", tariffID, err)
	}

	return &t, nil
}

func (s *TariffService) UpdateTariff(tariffID int, name string, price float64) error {
	_, err := s.db.Exec(
		"UPDATE lesbaza.tariffs SET name = $1, price_per_day = $2 WHERE tariff_id = $3",
		name, price, tariffID,
	)
	if err != nil {
		return fmt.Errorf("ошибка обновления тарифа: %w", err)
	}
	return nil
}

func (s *TariffService) DeleteTariff(tariffID int) error {
	// Проверяем, не используется ли тариф гостями
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM lesbaza.guests WHERE tariff_id = $1", tariffID).Scan(&count)
	if err != nil {
		return fmt.Errorf("ошибка проверки использования тарифа: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("нельзя удалить тариф: он используется %d гостями", count)
	}

	_, err = s.db.Exec("DELETE FROM lesbaza.tariffs WHERE tariff_id = $1", tariffID)
	if err != nil {
		return fmt.Errorf("ошибка удаления тарифа: %w", err)
	}
	return nil

}
