package ui

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/VallfIK/bazaotdx/internal/models"
	"github.com/VallfIK/bazaotdx/internal/service"
)

// BookingCalendar представляет календарный виджет для бронирования
type BookingCalendar struct {
	widget.BaseWidget

	bookingService *service.BookingService
	cottageService *service.CottageService
	tariffService  *service.TariffService

	currentMonth time.Time
	calendarData map[time.Time]map[int]models.BookingStatus
	cottages     []models.Cottage

	window    fyne.Window
	onRefresh func()

	// UI элементы
	monthLabel   *widget.Label
	calendarGrid *fyne.Container
}

// NewBookingCalendar создает новый календарь бронирования
func NewBookingCalendar(
	bookingService *service.BookingService,
	cottageService *service.CottageService,
	tariffService *service.TariffService,
	window fyne.Window,
) *BookingCalendar {
	bc := &BookingCalendar{
		bookingService: bookingService,
		cottageService: cottageService,
		tariffService:  tariffService,
		currentMonth:   time.Now().Local(),
		window:         window,
	}

	bc.ExtendBaseWidget(bc)
	bc.loadData()

	return bc
}

// loadData загружает данные для календаря
func (bc *BookingCalendar) loadData() {
	// Получаем список домиков
	cottages, err := bc.cottageService.GetAllCottages()
	if err != nil {
		dialog.ShowError(err, bc.window)
		return
	}
	bc.cottages = cottages

	// Получаем данные бронирования
	startDate := time.Date(bc.currentMonth.Year(), bc.currentMonth.Month(), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1)

	calendarData, err := bc.bookingService.GetCalendarData(startDate, endDate)
	if err != nil {
		dialog.ShowError(err, bc.window)
		return
	}
	bc.calendarData = calendarData
}

// CreateRenderer создает рендерер для календаря
func (bc *BookingCalendar) CreateRenderer() fyne.WidgetRenderer {
	// Заголовок с навигацией по месяцам
	bc.monthLabel = widget.NewLabel(bc.currentMonth.Format("January 2006"))
	// Устанавливаем русский формат для месяца
	bc.monthLabel.Text = time.Month(bc.currentMonth.Month()).String() + " " + bc.currentMonth.Format("2006")
	bc.monthLabel.TextStyle = fyne.TextStyle{Bold: true}
	bc.monthLabel.Alignment = fyne.TextAlignCenter
	bc.monthLabel.TextStyle = fyne.TextStyle{Bold: true}
	bc.monthLabel.Alignment = fyne.TextAlignCenter

	prevBtn := widget.NewButton("◀", func() {
		bc.currentMonth = bc.currentMonth.AddDate(0, -1, 0)
		bc.loadData()
		bc.updateCalendar()
	})

	nextBtn := widget.NewButton("▶", func() {
		bc.currentMonth = bc.currentMonth.AddDate(0, 1, 0)
		bc.loadData()
		bc.updateCalendar()
	})

	todayBtn := widget.NewButton("Сегодня", func() {
		bc.currentMonth = time.Now().Local()
		bc.loadData()
		bc.updateCalendar()
	})

	header := container.NewBorder(
		nil, nil,
		prevBtn,
		container.NewHBox(todayBtn, nextBtn),
		bc.monthLabel,
	)

	// Создаем сетку календаря
	bc.calendarGrid = bc.createCalendarGrid()

	// Легенда
	legend := bc.createLegend()

	// Основной контейнер
	content := container.NewBorder(
		container.NewVBox(header, legend),
		nil, nil, nil,
		container.NewScroll(bc.calendarGrid),
	)

	return widget.NewSimpleRenderer(content)
}

// createCalendarGrid создает сетку календаря
func (bc *BookingCalendar) createCalendarGrid() *fyne.Container {
	// Получаем первый день месяца в локальном часовом поясе
	startDate := time.Date(bc.currentMonth.Year(), bc.currentMonth.Month(), 1, 0, 0, 0, 0, time.Local)
	// Получаем последний день месяца
	endDate := startDate.AddDate(0, 1, -1)
	daysInMonth := endDate.Day()

	// Создаем сетку: первый столбец - названия домиков, остальные - дни
	gridObjects := []fyne.CanvasObject{}

	// Заголовок - пустая ячейка
	gridObjects = append(gridObjects, widget.NewLabel(""))

	// Заголовки дней
	for day := 1; day <= daysInMonth; day++ {
		date := time.Date(startDate.Year(), startDate.Month(), day, 0, 0, 0, 0, time.Local)
		dayLabel := widget.NewLabel(fmt.Sprintf("%d\n%s", day, bc.getWeekdayShort(date.Weekday())))
		dayLabel.Alignment = fyne.TextAlignCenter

		// Выделяем выходные
		if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
			dayLabel.TextStyle = fyne.TextStyle{Bold: true}
		}

		gridObjects = append(gridObjects, dayLabel)
	}

	// Строки для каждого домика
	for _, cottage := range bc.cottages {
		// Название домика
		cottageLabel := widget.NewLabel(cottage.Name)
		cottageLabel.TextStyle = fyne.TextStyle{Bold: true}
		gridObjects = append(gridObjects, cottageLabel)

		// Ячейки для каждого дня
		for day := 1; day <= daysInMonth; day++ {
			date := time.Date(startDate.Year(), startDate.Month(), day, 0, 0, 0, 0, time.Local)
			cell := bc.createCalendarCell(cottage.ID, date)
			gridObjects = append(gridObjects, cell)
		}
	}

	// Создаем сетку
	grid := container.New(layout.NewGridLayoutWithColumns(daysInMonth+1), gridObjects...)

	return grid
}

// createCalendarCell создает ячейку календаря
func (bc *BookingCalendar) createCalendarCell(cottageID int, date time.Time) fyne.CanvasObject {
	// Получаем статус бронирования для этого домика и дня
	var status models.BookingStatus
	if dayData, exists := bc.calendarData[date]; exists {
		if s, ok := dayData[cottageID]; ok {
			status = s
		}
	}

	// Определяем цвет ячейки
	bgColor := bc.getStatusColor(status)

	// Создаем прямоугольник фона
	bg := canvas.NewRectangle(bgColor)
	bg.SetMinSize(fyne.NewSize(80, 40))

	// Создаем контент ячейки
	var content fyne.CanvasObject

	if status.BookingID > 0 {
		// Есть бронирование
		text := ""
		if status.IsCheckIn {
			text = "→ " + bc.truncateString(status.GuestName, 8)
		} else if status.IsCheckOut {
			text = bc.truncateString(status.GuestName, 8) + " →"
		} else {
			// Для промежуточных дней показываем только имя
			text = bc.truncateString(status.GuestName, 10)
		}

		label := canvas.NewText(text, color.White)
		label.TextSize = 10
		label.Alignment = fyne.TextAlignCenter
		content = label
	} else {
		// Свободно - пустая ячейка
		content = widget.NewLabel("")
	}

	// Создаем кликабельный контейнер
	clickable := widget.NewButton("", func() {
		bc.onCellTapped(cottageID, date, status)
	})
	clickable.Importance = widget.LowImportance
	clickable.ExtendBaseWidget(clickable)

	// Объединяем все в один контейнер
	cell := container.NewStack(
		bg,
		container.NewCenter(content),
		clickable,
	)

	return cell
}

// getStatusColor возвращает цвет для статуса
func (bc *BookingCalendar) getStatusColor(status models.BookingStatus) color.Color {
	// Если нет бронирования, считаем день свободным
	if status.BookingID <= 0 {
		return color.NRGBA{R: 40, G: 167, B: 69, A: 255} // Зеленый (свободно)
	}

	switch status.Status {
	case models.BookingStatusBooked:
		return color.NRGBA{R: 255, G: 205, B: 83, A: 255} // Светло-желтый (бронь)
	case models.BookingStatusCheckedIn:
		return color.NRGBA{R: 0, G: 123, B: 255, A: 255} // Синий (заселен)
	case models.BookingStatusCancelled:
		return color.NRGBA{R: 220, G: 53, B: 69, A: 255} // Красный (отменено)
	default:
		return color.NRGBA{R: 255, G: 205, B: 83, A: 255} // Светло-желтый по умолчанию
	}
}

// onCellTapped обрабатывает клик по ячейке
func (bc *BookingCalendar) onCellTapped(cottageID int, date time.Time, status models.BookingStatus) {
	if status.BookingID > 0 {
		// Показываем информацию о брони
		bc.showBookingDetails(status.BookingID)
	} else {
		// Открываем форму быстрого бронирования
		bc.showQuickBookingForm(cottageID, date)
	}
}

// showBookingDetails показывает детали брони
func (bc *BookingCalendar) showBookingDetails(bookingID int) {
	booking, err := bc.bookingService.GetBookingByID(bookingID)
	if err != nil {
		dialog.ShowError(err, bc.window)
		return
	}

	// Находим название домика
	cottageName := ""
	for _, c := range bc.cottages {
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
				widget.NewLabel(fmt.Sprintf("Заезд: %s", booking.CheckInDate.Format("02.01.2006"))),
				widget.NewLabel(fmt.Sprintf("Выезд: %s", booking.CheckOutDate.Format("02.01.2006"))),
				widget.NewLabel(fmt.Sprintf("Статус: %s", bc.getStatusText(booking.Status))),
				widget.NewLabel(fmt.Sprintf("Стоимость: %.2f руб.", booking.TotalCost)),
			),
		),
	)

	// Кнопки действий в зависимости от статуса
	actions := container.NewHBox()

	switch booking.Status {
	case models.BookingStatusBooked:
		actions.Add(widget.NewButton("Заселить", func() {
			err := bc.bookingService.CheckInBooking(booking.ID)
			if err != nil {
				dialog.ShowError(err, bc.window)
				return
			}
			bc.Update()
			dialog.ShowInformation("Успешно", "Гость заселен", bc.window)
		}))
		actions.Add(widget.NewButton("Отменить", func() {
			dialog.ShowConfirm("Подтверждение", "Отменить бронирование?", func(ok bool) {
				if ok {
					err := bc.bookingService.CancelBooking(booking.ID)
					if err != nil {
						dialog.ShowError(err, bc.window)
						return
					}
					bc.Update()
					dialog.ShowInformation("Успешно", "Бронирование отменено", bc.window)
				}
			}, bc.window)
		}))
	}

	content.Add(actions)

	d := dialog.NewCustom("Детали бронирования", "Закрыть", content, bc.window)
	d.Resize(fyne.NewSize(400, 400))

	d.Show()
}

// showQuickBookingForm показывает форму быстрого бронирования
func (bc *BookingCalendar) showQuickBookingForm(cottageID int, startDate time.Time) {
	// Получаем информацию о домике
	cottages, err := bc.cottageService.GetAllCottages()
	if err != nil {
		dialog.ShowError(err, bc.window)
		return
	}
	var cottage models.Cottage
	for _, c := range cottages {
		if c.ID == cottageID {
			cottage = c
			break
		}
	}

	// Проверяем доступность домика на даты
	available, err := bc.bookingService.IsCottageAvailable(cottageID, startDate, startDate.AddDate(0, 0, 1))
	if err != nil {
		dialog.ShowError(err, bc.window)
		return
	}
	if !available {
		dialog.ShowError(fmt.Errorf("домик недоступен на выбранные даты"), bc.window)
		return
	}

	// Используем найденный домик из первого цикла
	for _, c := range bc.cottages {
		if c.ID == cottageID {
			cottage = c
			break
		}
	}

	// Получаем доступные тарифы
	tariffs, err := bc.tariffService.GetTariffs()
	if err != nil {
		dialog.ShowError(err, bc.window)
		return
	}

	// Создаем список тарифов для выпадающего списка
	tariffOptions := make([]string, len(tariffs))
	for i, tariff := range tariffs {
		tariffOptions[i] = fmt.Sprintf("%s - %.2f руб./сутки", tariff.Name, tariff.PricePerDay)
	}

	// Поля формы
	nameEntry := widget.NewEntry()
	phoneEntry := widget.NewEntry()
	emailEntry := widget.NewEntry()
	notesEntry := widget.NewMultiLineEntry()
	notesEntry.SetMinRowsVisible(3)

	// Выбор тарифа
	tariffSelect := widget.NewSelect(tariffOptions, nil)

	checkInDate := startDate
	checkOutDate := startDate.AddDate(0, 0, 1)

	// Расчет стоимости
	costLabel := widget.NewLabel("Стоимость: -")

	updateCost := func() {
		if tariffSelect.SelectedIndex() >= 0 {
			days := int(checkOutDate.Sub(checkInDate).Hours() / 24)
			if days <= 0 {
				days = 1
			}
			cost := float64(days) * tariffs[tariffSelect.SelectedIndex()].PricePerDay
			costLabel.SetText(fmt.Sprintf("Стоимость: %.2f руб. (%d дней)", cost, days))
		}
	}

	checkInPicker := NewDatePickerButton(
		"Дата заезда",
		bc.window,
		func(t time.Time) {
			checkInDate = t
			updateCost()
			// Проверяем доступность при изменении даты заезда
			if checkOutDate.After(checkInDate) {
				available, err := bc.bookingService.IsCottageAvailable(cottageID, checkInDate, checkOutDate)
				if err != nil {
					dialog.ShowError(err, bc.window)
					return
				}
				if !available {
					dialog.ShowError(fmt.Errorf("домик недоступен на выбранные даты"), bc.window)
					return
				}
			}
		},
	)
	checkInPicker.SetSelectedDate(checkInDate)

	checkOutPicker := NewDatePickerButton(
		"Дата выезда",
		bc.window,
		func(t time.Time) {
			checkOutDate = t
			updateCost()
			// Проверяем доступность при изменении даты выезда
			if checkOutDate.After(checkInDate) {
				available, err := bc.bookingService.IsCottageAvailable(cottageID, checkInDate, checkOutDate)
				if err != nil {
					dialog.ShowError(err, bc.window)
					return
				}
				if !available {
					dialog.ShowError(fmt.Errorf("домик недоступен на выбранные даты"), bc.window)
					return
				}
			}
		},
	)
	checkOutPicker.SetSelectedDate(checkOutDate)

	checkInPicker.onDateChange = func(t time.Time) {
		checkInDate = t
		updateCost()
	}
	checkOutPicker.onDateChange = func(t time.Time) {
		checkOutDate = t
		updateCost()
	}

	tariffSelect.OnChanged = func(_ string) {
		updateCost()
	}

	// Форма
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Домик", Widget: widget.NewLabel(cottage.Name)},
			{Text: "ФИО", Widget: nameEntry},
			{Text: "Телефон", Widget: phoneEntry},
			{Text: "Email", Widget: emailEntry},
			{Text: "Дата заезда", Widget: checkInPicker.button},
			{Text: "Дата выезда", Widget: checkOutPicker.button},
			{Text: "Тариф", Widget: tariffSelect},
			{Text: "", Widget: costLabel},
			{Text: "Примечания", Widget: notesEntry},
		},
		OnSubmit: func() {
			// Валидация
			if nameEntry.Text == "" || phoneEntry.Text == "" || tariffSelect.SelectedIndex() < 0 {
				dialog.ShowError(fmt.Errorf("Заполните обязательные поля"), bc.window)
				return
			}

			// Создаем бронь
			booking := models.Booking{
				CottageID:    cottageID,
				GuestName:    nameEntry.Text,
				Phone:        phoneEntry.Text,
				Email:        emailEntry.Text,
				CheckInDate:  checkInDate,
				CheckOutDate: checkOutDate,
				TariffID:     tariffs[tariffSelect.SelectedIndex()].ID,
				Notes:        notesEntry.Text,
			}

			// Рассчитываем стоимость
			days := int(checkOutDate.Sub(checkInDate).Hours() / 24)
			if days <= 0 {
				days = 1
			}
			booking.TotalCost = float64(days) * tariffs[tariffSelect.SelectedIndex()].PricePerDay

			// Сохраняем
			_, err := bc.bookingService.CreateBooking(booking)
			if err != nil {
				dialog.ShowError(err, bc.window)
				return
			}

			bc.Update()
			dialog.ShowInformation("Успешно", "Бронирование создано", bc.window)
		},
		OnCancel: func() {
			// Закрываем диалог
		},
	}

	d := dialog.NewCustom("Быстрое бронирование", "Закрыть", form, bc.window)
	d.Resize(fyne.NewSize(500, 600))
	d.Show()
}

// createLegend создает легенду
func (bc *BookingCalendar) createLegend() fyne.CanvasObject {
	items := []fyne.CanvasObject{
		widget.NewLabel("Легенда:"),
		bc.createLegendItem("Свободно", color.NRGBA{R: 40, G: 167, B: 69, A: 255}),
		bc.createLegendItem("Забронировано", color.NRGBA{R: 255, G: 193, B: 7, A: 255}),
		bc.createLegendItem("Заселено", color.NRGBA{R: 220, G: 53, B: 69, A: 255}),
		bc.createLegendItem("Выселено", color.NRGBA{R: 108, G: 117, B: 125, A: 255}),
	}

	return container.NewHBox(items...)
}

// createLegendItem создает элемент легенды
func (bc *BookingCalendar) createLegendItem(text string, color color.Color) fyne.CanvasObject {
	rect := canvas.NewRectangle(color)
	rect.SetMinSize(fyne.NewSize(20, 20))
	label := widget.NewLabel(text)
	return container.NewHBox(rect, label)
}

// getWeekdayShort возвращает короткое название дня недели
func (bc *BookingCalendar) getWeekdayShort(weekday time.Weekday) string {
	days := []string{"Вс", "Пн", "Вт", "Ср", "Чт", "Пт", "Сб"}
	return days[weekday]
}

// getStatusText возвращает текст статуса
func (bc *BookingCalendar) getStatusText(status string) string {
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

// truncateString обрезает строку до максимальной длины
func (bc *BookingCalendar) truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-2] + ".."
}

// updateCalendar обновляет отображение календаря
func (bc *BookingCalendar) updateCalendar() {
	// Обновляем текст месяца
	bc.monthLabel.Text = time.Month(bc.currentMonth.Month()).String() + " " + bc.currentMonth.Format("2006")
	// Обновляем сетку календаря
	bc.calendarGrid.Objects = bc.createCalendarGrid().Objects
}

// Update обновляет данные и отображение календаря
func (bc *BookingCalendar) Update() {
	bc.loadData()
	bc.updateCalendar()
	if bc.onRefresh != nil {
		bc.onRefresh()
	}
}

// SetOnRefresh устанавливает callback для обновления
func (bc *BookingCalendar) SetOnRefresh(f func()) {
	bc.onRefresh = f
}
