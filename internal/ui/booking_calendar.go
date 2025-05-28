package ui

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/VallfIK/bazaotdx/internal/models"
	"github.com/VallfIK/bazaotdx/internal/service"
)

var (
	GreenColor = color.NRGBA{R: 40, G: 167, B: 69, A: 255}
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

// SetOnRefresh устанавливает callback для обновления
func (bc *BookingCalendar) SetOnRefresh(f func()) {
	bc.onRefresh = f
}

// Update обновляет данные и отображение календаря
func (bc *BookingCalendar) Update() {
	bc.loadData()
	bc.updateCalendar()
	if bc.onRefresh != nil {
		bc.onRefresh()
	}
}

// ClickableRect - прямоугольник который можно кликать без hover эффекта
type ClickableRect struct {
	widget.BaseWidget
	rect     *canvas.Rectangle
	onTapped func()
}

func NewClickableRect(bgColor color.Color, onTapped func()) *ClickableRect {
	r := &ClickableRect{
		rect:     canvas.NewRectangle(bgColor),
		onTapped: onTapped,
	}
	r.ExtendBaseWidget(r)
	return r
}

func (r *ClickableRect) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(r.rect)
}

func (r *ClickableRect) Tapped(_ *fyne.PointEvent) {
	if r.onTapped != nil {
		r.onTapped()
	}
}

func (r *ClickableRect) TappedSecondary(_ *fyne.PointEvent) {}

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
func (bc *BookingCalendar) loadData() error {
	// Получаем список домиков
	cottages, err := bc.cottageService.GetAllCottages()
	if err != nil {
		return fmt.Errorf("error loading cottages: %v", err)
	}
	bc.cottages = cottages

	// Получаем данные бронирования для всего месяца
	startDate := time.Date(bc.currentMonth.Year(), bc.currentMonth.Month(), 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(bc.currentMonth.Year(), bc.currentMonth.Month()+1, 1, 0, 0, 0, 0, time.Local).AddDate(0, 0, -1)

	calendarData, err := bc.bookingService.GetCalendarData(startDate, endDate)
	if err != nil {
		return fmt.Errorf("error loading calendar data: %v", err)
	}
	bc.calendarData = calendarData

	return nil
}

// CreateRenderer создает рендерер для календаря
func (bc *BookingCalendar) CreateRenderer() fyne.WidgetRenderer {
	// Заголовок с навигацией по месяцам
	bc.monthLabel = widget.NewLabel("")
	bc.updateMonthLabel()
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

// updateMonthLabel обновляет текст месяца на русском
func (bc *BookingCalendar) updateMonthLabel() {
	months := []string{
		"", "Январь", "Февраль", "Март", "Апрель", "Май", "Июнь",
		"Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь",
	}
	month := months[bc.currentMonth.Month()]
	year := bc.currentMonth.Year()
	bc.monthLabel.SetText(fmt.Sprintf("%s %d", month, year))
}

// createCalendarGrid создает сетку календаря
func (bc *BookingCalendar) createCalendarGrid() *fyne.Container {
	today := time.Now().Local().Truncate(24 * time.Hour)
	firstDay := time.Date(bc.currentMonth.Year(), bc.currentMonth.Month(), 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, -1)

	// Определяем диапазон дней для отображения
	startDay := 1
	if firstDay.Before(today) {
		startDay = today.Day()
	}

	daysToShow := lastDay.Day() - startDay + 1

	// Создаем основной контейнер
	mainContainer := container.NewVBox()

	// Создаем заголовок с днями недели
	daysHeader := container.NewGridWithColumns(daysToShow + 1)

	// Первая ячейка - заголовок "Домик"
	daysHeader.Add(widget.NewCard("", "", widget.NewLabel("Домик")))

	// Добавляем заголовки дней (только текущие и будущие)
	for day := startDay; day <= lastDay.Day(); day++ {
		currentDate := time.Date(bc.currentMonth.Year(), bc.currentMonth.Month(), day, 0, 0, 0, 0, time.Local)
		weekday := bc.getWeekdayShort(currentDate.Weekday())

		dayCard := widget.NewCard("", "",
			container.NewVBox(
				widget.NewLabel(fmt.Sprintf("%d", day)),
				widget.NewLabel(weekday),
			),
		)
		daysHeader.Add(dayCard)
	}

	mainContainer.Add(daysHeader)

	// Создаем строки для каждого домика
	for _, cottage := range bc.cottages {
		cottageRow := bc.createCottageRow(cottage, startDay, lastDay.Day())
		mainContainer.Add(cottageRow)
	}

	return mainContainer
}

// createCottageRow создает строку для домика
func (bc *BookingCalendar) createCottageRow(cottage models.Cottage, startDay, endDay int) *fyne.Container {
	daysToShow := endDay - startDay + 1

	// Создаем сетку для строки
	row := container.NewGridWithColumns(daysToShow + 1)

	// Первая ячейка - название домика
	cottageNameCard := widget.NewCard("", "", widget.NewLabel(cottage.Name))
	cottageNameCard.Resize(fyne.NewSize(100, 60))
	row.Add(cottageNameCard)

	// Добавляем ячейки для каждого дня
	for day := startDay; day <= endDay; day++ {
		currentDate := time.Date(bc.currentMonth.Year(), bc.currentMonth.Month(), day, 0, 0, 0, 0, time.Local)
		cell := bc.createCalendarCell(cottage.ID, currentDate)
		row.Add(cell)
	}

	return row
}

// createCalendarCell создает ячейку календаря
func (bc *BookingCalendar) createCalendarCell(cottageID int, date time.Time) fyne.CanvasObject {
	// Получаем статус для этого домика и дня
	var status *models.BookingStatus

	if dayData, exists := bc.calendarData[date]; exists {
		if cottageStatus, ok := dayData[cottageID]; ok {
			status = &cottageStatus
		}
	}

	// Создаем кнопку
	if status != nil && status.IsCheckOut {
		// День выезда - диагональная кнопка
		return bc.createDiagonalCell(cottageID, date, *status)
	} else {
		// Обычная кнопка
		return bc.createRegularCell(cottageID, date, status)
	}
}

// createRegularCell создает обычную ячейку
func (bc *BookingCalendar) createRegularCell(cottageID int, date time.Time, status *models.BookingStatus) fyne.CanvasObject {
	var bgColor color.Color
	var text string

	if status != nil && status.BookingID > 0 {
		bgColor = bc.getStatusColor(*status)
		if status.IsCheckIn {
			text = "→ " + bc.truncateString(status.GuestName, 6)
		} else {
			text = bc.truncateString(status.GuestName, 8)
		}
	} else {
		// Свободный день
		bgColor = color.NRGBA{R: 40, G: 167, B: 69, A: 255} // Зеленый
		text = "+"
	}

	// Создаем кликабельный прямоугольник
	clickableRect := NewClickableRect(bgColor, func() {
		if status != nil {
			bc.onCellTapped(cottageID, date, *status)
		} else {
			bc.onCellTapped(cottageID, date, models.BookingStatus{})
		}
	})

	// Добавляем текст
	label := canvas.NewText(text, color.White)
	label.TextSize = 10
	label.Alignment = fyne.TextAlignCenter

	content := container.NewStack(
		clickableRect,
		container.NewCenter(label),
	)

	content.Resize(fyne.NewSize(35, 60))
	return content
}

// createDiagonalCell создает диагональную ячейку для дня выезда
func (bc *BookingCalendar) createDiagonalCell(cottageID int, date time.Time, status models.BookingStatus) fyne.CanvasObject {
	var leftColor, rightColor color.Color
	var text string

	// День выезда - левая часть показывает текущую бронь, правая - свободна для новой брони
	leftColor = bc.getStatusColor(status)  // Цвет текущей брони (до 12:00)
	rightColor = bc.getStatusColor(status) // Цвет текущей брони (после 12:00)
	text = bc.truncateString(status.GuestName, 6) + " → до 12:00"

	button := NewDiagonalButton(leftColor, rightColor, text,
		func() {
			// Левая часть - показываем детали текущей брони
			bc.onCellTapped(cottageID, date, status)
		},
		func() {
			// Правая часть - открываем форму бронирования с этого дня
			bc.showQuickBookingForm(cottageID, date)
		})

	button.Resize(fyne.NewSize(35, 60))
	return button
}

// getStatusColor возвращает цвет для статуса
func (bc *BookingCalendar) getStatusColor(status models.BookingStatus) color.Color {
	switch status.Status {
	case models.BookingStatusBooked:
		return color.NRGBA{R: 255, G: 193, B: 7, A: 255} // Желтый (бронь)
	case models.BookingStatusCheckedIn:
		return color.NRGBA{R: 0, G: 123, B: 255, A: 255} // Синий (заселен)
	case models.BookingStatusCheckedOut:
		return color.NRGBA{R: 108, G: 117, B: 125, A: 255} // Серый (выселен)
	case models.BookingStatusCancelled:
		return color.NRGBA{R: 220, G: 53, B: 69, A: 255} // Красный (отменено)
	default:
		return color.NRGBA{R: 40, G: 167, B: 69, A: 255} // Зеленый (свободно)
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
		dialog.ShowError(fmt.Errorf("error loading booking details: %v", err), bc.window)
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
	case models.BookingStatusCheckedIn:
		actions.Add(widget.NewButton("Выселить", func() {
			err := bc.bookingService.CheckOutBooking(booking.ID)
			if err != nil {
				dialog.ShowError(err, bc.window)
				return
			}
			bc.Update()
			dialog.ShowInformation("Успешно", "Гость выселен", bc.window)
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
	var cottage models.Cottage
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

	checkInDate := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 14, 0, 0, 0, time.Local)
	checkOutDate := checkInDate.AddDate(0, 0, 1)
	checkOutDate = time.Date(checkOutDate.Year(), checkOutDate.Month(), checkOutDate.Day(), 12, 0, 0, 0, time.Local)

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
		},
	)
	checkInPicker.SetSelectedDate(checkInDate)

	checkOutPicker := NewDatePickerButton(
		"Дата выезда",
		bc.window,
		func(t time.Time) {
			checkOutDate = t
			updateCost()
		},
	)
	checkOutPicker.SetSelectedDate(checkOutDate)

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

			if checkOutDate.Before(checkInDate) || checkOutDate.Equal(checkInDate) {
				dialog.ShowError(fmt.Errorf("Дата выезда должна быть позже даты заезда"), bc.window)
				return
			}

			// Проверяем доступность
			available, err := bc.bookingService.IsCottageAvailable(cottageID, checkInDate, checkOutDate)
			if err != nil {
				dialog.ShowError(err, bc.window)
				return
			}
			if !available {
				dialog.ShowError(fmt.Errorf("домик недоступен на выбранные даты"), bc.window)
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
			_, err = bc.bookingService.CreateBooking(booking)
			if err != nil {
				dialog.ShowError(err, bc.window)
				return
			}

			bc.Update()
			dialog.ShowInformation("Успешно", "Бронирование создано", bc.window)
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
		bc.createLegendItem("Заселено", color.NRGBA{R: 0, G: 123, B: 255, A: 255}),
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
	bc.updateMonthLabel()
	bc.calendarGrid.Objects = bc.createCalendarGrid().Objects
	bc.calendarGrid.Refresh()
}
