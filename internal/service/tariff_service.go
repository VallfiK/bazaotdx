package service

import (
	"database/sql"

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
	return err
}

func (s *TariffService) GetTariffs() ([]models.Tariff, error) {
	rows, err := s.db.Query("SELECT tariff_id, name, price_per_day FROM lesbaza.tariffs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tariffs []models.Tariff
	for rows.Next() {
		var t models.Tariff
		err := rows.Scan(&t.ID, &t.Name, &t.PricePerDay)
		if err != nil {
			return nil, err
		}
		tariffs = append(tariffs, t)
	}
	return tariffs, nil
}
