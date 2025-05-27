package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/VallfIK/bazaotdx/internal/models"
	"github.com/VallfIK/bazaotdx/internal/service"
)

// BookingListWidget виджет для отображения списка бронирований
type BookingListWidget struct {
	widget.BaseWidget

	bookingService *service.BookingService
	cottageService *service.CottageService

	window   fyne.Window
	bookings []models.Booking
	cottages []models.Cottage

	// UI элементы
	list         *widget.List
	filterSelect *widget.Select
	refreshBtn   *widget.Button

	onRefresh func()
}

// NewBookingListWidget создает новый виджет списка бронирований
func NewBookingListWidget(
	bookingService *service.BookingService,
	cottageService *service.CottageService,
	window fyne.Window,
) *BookingListWidget {
	blw := &BookingListWidget{
		bookingService: bookingService,
		cottageService: cottageService,
		window:         window,
	}

	blw.ExtendBaseWidget(blw)
	blw.loadData()

	return blw
}

// loadData загружает данные
func (blw *BookingListWidget) loadData() {
	// Загружаем домики
	cottages, err := blw.cottageService.GetAllCottages()
	if err != nil {
		dialog.ShowError(err, blw.window)
		return
	}
	blw.cottages = cottages

	// Загружаем брони в зависимости от фильтра
	blw.loadBookings()
}

// loadBookings загружает брони в зависимости от выбранного фильтра
func (blw *BookingListWidget) loadBookings() {
	var bookings []models.Booking
	var err error

	filter := "Предстоящие"
	if blw.filterSelect != nil && blw.filterSelect.Selected != "" {
		filter = blw.filterSelect.Selected
	}

	today := time.Now().Truncate(24 * time.Hour)

	switch filter {
	case "Предстоящие":
		bookings, err = blw.bookingService.GetBookingsByStatus(models.BookingStatusBooked, today)
	case "Текущие":
		bookings, err = blw.bookingService.GetBookingsByStatus(models.BookingStatusCheckedIn, time.Time{})
	case "Все активные":
		// Получаем все брони за следующие 3 месяца
		endDate := today.AddDate(0, 3, 0)
		bookings, err = blw.bookingService.GetBookingsByDateRange(today, endDate)
	case "За последний месяц":
		startDate := today.AddDate(0, -1, 0)
		bookings, err = blw.bookingService.GetBookingsByDateRange(startDate, today)
	default:
		bookings = []models.Booking{}
	}

	if err != nil {
		dialog.ShowError(err, blw.window)
		return
	}

	blw.bookings = bookings
	if blw.list != nil {
		blw.list.Refresh()
	}
}

// CreateRenderer создает рендерер для виджета
func (blw *BookingListWidget) CreateRenderer() fyne.WidgetRenderer {
	// Фильтр
	blw.filterSelect = widget.NewSelect(
		[]string{"Предстоящие", "Текущие", "Все активные", "За последний месяц"},
		func(value string) {
			blw.loadBookings()
		},
	)
	blw.filterSelect.SetSelected("Предстоящие")

	// Кнопка обновления
	blw.refreshBtn = widget.NewButton("Обновить", func() {
		blw.loadData()
		dialog.ShowInformation("Успешно", "Данные обновлены", blw.window)
	})

	// Заголовок
	header := container.NewBorder(
		nil, nil,
		widget.NewLabel("Фильтр:"),
		blw.refreshBtn,
		blw.filterSelect,
	)

	// Список бронирований
	blw.list = widget.NewList(
		func() int {
			return len(blw.bookings)
		},
		func() fyne.CanvasObject {
			return blw.createBookingItem()
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			blw.updateBookingItem(id, item)
		},
	)

	blw.list.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(blw.bookings) {
			blw.showBookingActions(blw.bookings[id])
		}
		blw.list.Unselect(id)
	}

	content := container.NewBorder(header, nil, nil, nil, blw.list)

	return widget.NewSimpleRenderer(content)
}

// createBookingItem создает элемент списка
func (blw *BookingListWidget) createBookingItem() fyne.CanvasObject {
	nameLabel := widget.NewLabel("")
	nameLabel.TextStyle = fyne.TextStyle{Bold: true}

	cottageLabel := widget.NewLabel("")
	datesLabel := widget.NewLabel("")
	phoneLabel := widget.NewLabel("")
	statusLabel := widget.NewLabel("")

	info := container.NewVBox(
		nameLabel,
		container.NewHBox(cottageLabel, widget.NewLabel("-"), datesLabel),
		container.NewHBox(phoneLabel, widget.NewLabel("-"), statusLabel),
	)

	return widget.NewCard("", "", info)
}

// updateBookingItem обновляет данные элемента списка
func (blw *BookingListWidget) updateBookingItem(id widget.ListItemID, item fyne.CanvasObject) {
	if id >= len(blw.bookings) {
		return
	}

	booking := blw.bookings[id]
	card := item.(*widget.Card)
	content := card.Content.(*fyne.Container)

	// Находим название домика
	cottageName := ""
	for _, c := range blw.cottages {
		if c.ID == booking.CottageID {
			cottageName = c.Name
			break
		}
	}

	// Обновляем labels
	vbox := content.Objects[0].(*fyne.Container)

	nameLabel := vbox.Objects[0].(*widget.Label)
	nameLabel.SetText(booking.GuestName)

	infoBox1 := vbox.Objects[1].(*fyne.Container)
	cottageLabel := infoBox1.Objects[0].(*widget.Label)
	cottageLabel.SetText(cottageName)

	datesLabel := infoBox1.Objects[2].(*widget.Label)
	datesLabel.SetText(fmt.Sprintf("%s - %s",
		booking.CheckInDate.Format("02.01"),
		booking.CheckOutDate.Format("02.01"),
	))

	infoBox2 := vbox.Objects[2].(*fyne.Container)
	phoneLabel := infoBox2.Objects[0].(*widget.Label)
	phoneLabel.SetText(booking.Phone)

	statusLabel := infoBox2.Objects[2].(*widget.Label)
	statusLabel.SetText(blw.getStatusText(booking.Status))

	// Цвет карточки в зависимости от статуса
	switch booking.Status {
	case models.BookingStatusBooked:
		card.SetSubTitle("Забронировано")
	case models.BookingStatusCheckedIn:
		card.SetSubTitle("Заселено")
	case models.BookingStatusCheckedOut:
		card.SetSubTitle("Выселено")
	case models.BookingStatusCancelled:
		card.SetSubTitle("Отменено")
	}
}

// showBookingActions показывает действия для брони
func (blw *BookingListWidget) showBookingActions(booking models.Booking) {
	// Находим название домика
	cottageName := ""
	for _, c := range blw.cottages {
		if c.ID == booking.CottageID {
			cottageName = c.Name
			break
		}
	}

	content := container.NewVBox(
		widget.NewCard("Информация о брони", "",
			container.NewVBox(
				widget.NewLabel(fmt.Sprintf("Гость: %s", booking.GuestName)),
				widget.NewLabel(fmt.Sprintf("Телефон: %s", booking.Phone)),
				widget.NewLabel(fmt.Sprintf("Email: %s", booking.Email)),
				widget.NewLabel(fmt.Sprintf("Домик: %s", cottageName)),
				widget.NewLabel(fmt.Sprintf("Заезд: %s", booking.CheckInDate.Format("02.01.2006 15:04"))),
				widget.NewLabel(fmt.Sprintf("Выезд: %s", booking.CheckOutDate.Format("02.01.2006 15:04"))),
				widget.NewLabel(fmt.Sprintf("Статус: %s", blw.getStatusText(booking.Status))),
				widget.NewLabel(fmt.Sprintf("Стоимость: %.2f руб.", booking.TotalCost)),
				widget.NewLabel(fmt.Sprintf("Создано: %s", booking.CreatedAt.Format("02.01.2006 15:04"))),
			),
		),
	)

	if booking.Notes != "" {
		content.Add(widget.NewCard("", "Примечания",
			widget.NewLabel(booking.Notes),
		))
	}

	// Кнопки действий
	actions := container.NewHBox()

	switch booking.Status {
	case models.BookingStatusBooked:
		actions.Add(widget.NewButton("Заселить", func() {
			dialog.ShowConfirm("Подтверждение", "Заселить гостя?", func(ok bool) {
				if ok {
					err := blw.bookingService.CheckInBooking(booking.ID)
					if err != nil {
						dialog.ShowError(err, blw.window)
						return
					}
					blw.loadData()
					blw.triggerRefresh()
					dialog.ShowInformation("Успешно", "Гость заселен", blw.window)
				}
			}, blw.window)
		}))

		actions.Add(widget.NewButton("Отменить", func() {
			dialog.ShowConfirm("Подтверждение", "Отменить бронирование?", func(ok bool) {
				if ok {
					err := blw.bookingService.CancelBooking(booking.ID)
					if err != nil {
						dialog.ShowError(err, blw.window)
						return
					}
					blw.loadData()
					blw.triggerRefresh()
					dialog.ShowInformation("Успешно", "Бронирование отменено", blw.window)
				}
			}, blw.window)
		}))

		// Здесь можно добавить кнопку "Изменить даты"
	}

	content.Add(actions)

	d := dialog.NewCustom("Действия с бронированием", "Закрыть", content, blw.window)
	d.Resize(fyne.NewSize(500, 500))
	d.Show()
}

// getStatusText возвращает текст статуса
func (blw *BookingListWidget) getStatusText(status string) string {
	switch status {
	case models.BookingStatusBooked:
		return "Забронировано"
	case models.BookingStatusCheckedIn:
		return "Заселено"
	case models.BookingStatusCheckedOut:
		return "Выселено"
	case models.BookingStatusCancelled:
		return "Отменено"
	default:
		return status
	}
}

// SetOnRefresh устанавливает callback для обновления
func (blw *BookingListWidget) SetOnRefresh(f func()) {
	blw.onRefresh = f
}

// triggerRefresh вызывает callback обновления
func (blw *BookingListWidget) triggerRefresh() {
	if blw.onRefresh != nil {
		blw.onRefresh()
	}
}
