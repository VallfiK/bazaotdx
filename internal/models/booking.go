package models

import "time"

type Booking struct {
	ID           int       `db:"booking_id"`
	CottageID    int       `db:"cottage_id"`
	GuestName    string    `db:"guest_name"`
	Phone        string    `db:"phone"`
	Email        string    `db:"email"`
	CheckInDate  time.Time `db:"check_in_date"`
	CheckOutDate time.Time `db:"check_out_date"`
	Status       string    `db:"status"` // 'booked', 'checked_in', 'checked_out', 'cancelled'
	CreatedAt    time.Time `db:"created_at"`
	Notes        string    `db:"notes"`
	TariffID     int       `db:"tariff_id"`
	TotalCost    float64   `db:"total_cost"`
}

// BookingStatus константы для статусов
const (
	BookingStatusBooked     = "booked"
	BookingStatusCheckedIn  = "checked_in"
	BookingStatusCheckedOut = "checked_out"
	BookingStatusCancelled  = "cancelled"
	BookingStatusTemporary  = "temporary"
)

// CalendarDay представляет день в календаре
type CalendarDay struct {
	Date     time.Time
	Cottages map[int]BookingStatus // cottage_id -> status
}

// BookingStatus для отображения в календаре
type BookingStatus struct {
	Status     string
	BookingID  int
	GuestName  string
	IsPartDay  bool // true если это день заезда или выезда
	IsCheckIn  bool // true если это день заезда
	IsCheckOut bool // true если это день выезда
}
