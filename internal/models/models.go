package models

import "time"

type Cottage struct {
	ID     int
	Name   string
	Status string
}

type Guest struct {
	ID               int       `db:"guest_id"`
	FullName         string    `db:"full_name"`
	Email            string    `db:"email"`
	Phone            string    `db:"phone"`
	CottageID        int       `db:"cottage_id"`
	DocumentScanPath string    `db:"document_scan_path"`
	CheckInDate      time.Time `db:"check_in_date"`
	CheckOutDate     time.Time `db:"check_out_date"`
	TariffID         int       `db:"tariff_id"`
}

type Tariff struct {
	ID          int     `db:"tariff_id"`
	Name        string  `db:"name"`
	PricePerDay float64 `db:"price_per_day"`
}
