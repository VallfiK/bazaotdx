package models

type Cottage struct {
	ID     int
	Name   string
	Status string
}

type Guest struct {
	ID               int    `db:"guest_id"`
	FullName         string `db:"full_name"`
	Email            string `db:"email"`
	Phone            string `db:"phone"`
	CottageID        int    `db:"cottage_id"`
	DocumentScanPath string `db:"document_scan_path"`
}
