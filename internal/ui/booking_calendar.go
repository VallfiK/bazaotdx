package ui

import (
	"fmt"
	"image/color"
	"sort"
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

	// Сортируем домики по ID
	cottages := make([]models.Cottage, len(bc.cottages))
	copy(cottages, bc.cottages)
	sort.Slice(cottages, func(i, j int) bool {
		return cottages[i].ID < cottages[j].ID
	})

	// Создаем строки для каждого домика
	for _, cottage := range cottages {
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

	// Проверяем, нужна ли диагональная кнопка
	needsDiagonal := bc.checkIfNeedsDiagonal(cottageID, date)

	if needsDiagonal {
		return bc.createDiagonalCell(cottageID, date, status)
	} else {
		return bc.createRegularCell(cottageID, date, status)
	}
}

// checkIfNeedsDiagonal проверяет, нужна ли диагональная кнопка для этого дня
func (bc *BookingCalendar) checkIfNeedsDiagonal(cottageID int, date time.Time) bool {
	// Получаем все брони для этого домика
	startDate := date.AddDate(0, 0, -1) // Проверяем день до
	endDate := date.AddDate(0, 0, 1)    // И день после

	bookings, err := bc.bookingService.GetBookingsByDateRange(startDate, endDate)
	if err != nil {
		return false
	}

	// Проверяем, есть ли бронь которая заканчивается в этот день
	hasCheckOut := false
	// И может ли быть новая бронь в этот же день
	canCheckIn := true

	for _, booking := range bookings {
		if booking.CottageID != cottageID {
			continue
		}

		bookingCheckOut := time.Date(booking.CheckOutDate.Year(), booking.CheckOutDate.Month(), booking.CheckOutDate.Day(), 0, 0, 0, 0, time.Local)
		bookingCheckIn := time.Date(booking.CheckInDate.Year(), booking.CheckInDate.Month(), booking.CheckInDate.Day(), 0, 0, 0, 0, time.Local)

		// Если бронь заканчивается в этот день
		if bookingCheckOut.Equal(date) && booking.Status != models.BookingStatusCancelled {
			hasCheckOut = true
		}

		// Если есть бронь которая покрывает этот день полностью
		if bookingCheckIn.Equal(date) && booking.Status != models.BookingStatusCancelled {
			canCheckIn = false
		}
	}

	// Диагональная кнопка нужна, если есть выезд И можно заселиться
	return hasCheckOut && canCheckIn
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
		bgColor = GreenColor
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

// createDiagonalCell создает диагональную ячейку для дня выезда/заезда
func (bc *BookingCalendar) createDiagonalCell(cottageID int, date time.Time, status *models.BookingStatus) fyne.CanvasObject {
	var topColor, bottomColor color.Color
	var text string

	// Находим бронь, которая заканчивается в этот день
	checkoutBooking := bc.findCheckoutBooking(cottageID, date)

	// Проверяем, есть ли бронь на заезд в этот же день
	checkinBooking := bc.findCheckinBooking(cottageID, date)

	// Определяем цвета для диагональной кнопки
	if checkoutBooking != nil {
		// Есть существующая бронь (выезд до 12:00)
		bottomColor = bc.getStatusColorForBooking(*checkoutBooking)
	} else {
		// Нет брони на выезд
		bottomColor = color.NRGBA{R: 200, G: 200, B: 200, A: 255} // Серый
	}

	// Верхняя часть - для заезда после 14:00
	if checkinBooking != nil {
		// Есть бронь на заезд в этот день
		topColor = bc.getStatusColorForBooking(*checkinBooking)
	} else if status != nil && status.BookingID > 0 && status.IsCheckIn {
		// Используем статус из календарных данных если это день заезда
		topColor = bc.getStatusColor(*status)
	} else {
		// Свободно для заезда
		topColor = GreenColor
	}

	text = "←12:00\n14:00→"

	button := NewDiagonalButton(topColor, bottomColor, text,
		func() {
			// Верхняя часть - клик на заезд (после 14:00)
			if checkinBooking != nil {
				// Если есть бронь на заезд, показываем детали
				bc.showBookingDetails(checkinBooking.ID)
			} else if status != nil && status.BookingID > 0 && status.IsCheckIn {
				// Если есть статус заезда, показываем детали
				bc.showBookingDetails(status.BookingID)
			} else {
				// Открываем форму бронирования с этого дня (после 14:00)
				checkInTime := time.Date(date.Year(), date.Month(), date.Day(), 14, 0, 0, 0, time.Local)
				bc.showQuickBookingForm(cottageID, checkInTime)
			}
		},
		func() {
			// Нижняя часть - клик на выезд (до 12:00)
			if checkoutBooking != nil {
				bc.showBookingDetails(checkoutBooking.ID)
			} else if status != nil && status.BookingID > 0 {
				bc.onCellTapped(cottageID, date, *status)
			} else {
				dialog.ShowInformation("Информация",
					"В этот день происходит выезд из домика до 12:00.\nВы можете забронировать заезд с 14:00.",
					bc.window)
			}
		})

	button.Resize(fyne.NewSize(35, 60))
	return button
}

// findCheckinBooking находит бронь, которая начинается в указанный день
func (bc *BookingCalendar) findCheckinBooking(cottageID int, date time.Time) *models.Booking {
	startDate := date
	endDate := date.AddDate(0, 0, 1) // Следующий день

	bookings, err := bc.bookingService.GetBookingsByDateRange(startDate, endDate)
	if err != nil {
		return nil
	}

	for _, booking := range bookings {
		if booking.CottageID != cottageID || booking.Status == models.BookingStatusCancelled {
			continue
		}

		checkInDate := time.Date(booking.CheckInDate.Year(), booking.CheckInDate.Month(), booking.CheckInDate.Day(), 0, 0, 0, 0, time.Local)

		if checkInDate.Equal(date) {
			return &booking
		}
	}

	return nil
}

// findCheckoutBooking находит бронь, которая заканчивается в указанный день
func (bc *BookingCalendar) findCheckoutBooking(cottageID int, date time.Time) *models.Booking {
	startDate := date.AddDate(0, 0, -7) // Ищем за неделю до
	endDate := date.AddDate(0, 0, 1)    // И день после

	bookings, err := bc.bookingService.GetBookingsByDateRange(startDate, endDate)
	if err != nil {
		return nil
	}

	for _, booking := range bookings {
		if booking.CottageID != cottageID || booking.Status == models.BookingStatusCancelled {
			continue
		}

		checkOutDate := time.Date(booking.CheckOutDate.Year(), booking.CheckOutDate.Month(), booking.CheckOutDate.Day(), 0, 0, 0, 0, time.Local)

		if checkOutDate.Equal(date) {
			return &booking
		}
	}

	return nil
}

// getStatusColorForBooking возвращает цвет для брони
func (bc *BookingCalendar) getStatusColorForBooking(booking models.Booking) color.Color {
	switch booking.Status {
	case models.BookingStatusBooked:
		return color.NRGBA{R: 255, G: 193, B: 7, A: 255} // Желтый (бронь)
	case models.BookingStatusCheckedIn:
		return color.NRGBA{R: 255, G: 102, B: 0, A: 255} // Оранжевый (заселен)
	case models.BookingStatusCancelled:
		return color.NRGBA{R: 220, G: 53, B: 69, A: 255} // Красный (отменено)
	case models.BookingStatusCompleted:
		return GreenColor // Зеленый (завершено = свободно)
	default:
		return GreenColor // Зеленый (свободно)
	}
}

// getStatusColor возвращает цвет для статуса
func (bc *BookingCalendar) getStatusColor(status models.BookingStatus) color.Color {
	// Если это день заезда, используем цвет бронирования
	if status.IsCheckIn {
		return color.NRGBA{R: 255, G: 193, B: 7, A: 255} // Желтый (бронь)
	}

	switch status.Status {
	case models.BookingStatusBooked:
		return color.NRGBA{R: 255, G: 193, B: 7, A: 255} // Желтый (бронь)
	case models.BookingStatusCheckedIn:
		return color.NRGBA{R: 255, G: 102, B: 0, A: 255} // Оранжевый (заселен)
	case models.BookingStatusCancelled:
		return color.NRGBA{R: 220, G: 53, B: 69, A: 255} // Красный (отменено)
	case models.BookingStatusCompleted:
		return GreenColor // Зеленый (завершено = свободно)
	default:
		return GreenColor // Зеленый (свободно)
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
			dialog.ShowConfirm("Подтверждение",
				fmt.Sprintf("Выселить гостя %s из домика %s?", booking.GuestName, cottageName),
				func(ok bool) {
					if ok {
						err := bc.bookingService.CheckOutBooking(booking.ID)
						if err != nil {
							dialog.ShowError(err, bc.window)
							return
						}
						bc.Update()
						dialog.ShowInformation("Успешно", "Гость выселен", bc.window)
					}
				}, bc.window)
		}))

		// Кнопка для раннего выселения
		actions.Add(widget.NewButton("Ранний выезд", func() {
			bc.showEarlyCheckoutDialog(booking, cottageName)
		}))
	}

	content.Add(actions)

	d := dialog.NewCustom("Детали бронирования", "Закрыть", content, bc.window)
	d.Resize(fyne.NewSize(400, 400))
	d.Show()
}

// showQuickBookingForm показывает форму быстрого бронирования
func (bc *BookingCalendar) showQuickBookingForm(cottageID int, startDateTime time.Time) {
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

	if len(tariffs) == 0 {
		dialog.ShowError(fmt.Errorf("нет доступных тарифов"), bc.window)
		return
	}

	// Создаем список тарифов для выпадающего списка
	tariffOptions := make([]string, len(tariffs))
	for i, tariff := range tariffs {
		tariffOptions[i] = fmt.Sprintf("%s - %.2f руб./сутки", tariff.Name, tariff.PricePerDay)
	}

	// Поля формы
	nameEntry := widget.NewEntry()
	nameEntry.PlaceHolder = "ФИО гостя"

	phoneEntry := widget.NewEntry()
	phoneEntry.PlaceHolder = "+7 (999) 123-45-67"

	emailEntry := widget.NewEntry()
	emailEntry.PlaceHolder = "email@example.com"

	notesEntry := widget.NewMultiLineEntry()
	notesEntry.SetMinRowsVisible(3)
	notesEntry.PlaceHolder = "Дополнительные примечания..."

	// Выбор тарифа
	tariffSelect := widget.NewSelect(tariffOptions, nil)
	tariffSelect.SetSelected(tariffOptions[0]) // Выбираем первый тариф по умолчанию

	// Устанавливаем время заезда
	var checkInDate time.Time
	if startDateTime.Hour() == 0 && startDateTime.Minute() == 0 {
		// Если передана только дата, устанавливаем время 14:00
		checkInDate = time.Date(startDateTime.Year(), startDateTime.Month(), startDateTime.Day(), 14, 0, 0, 0, time.Local)
	} else {
		// Используем переданное время
		checkInDate = startDateTime
	}

	checkOutDate := checkInDate.AddDate(0, 0, 1)
	checkOutDate = time.Date(checkOutDate.Year(), checkOutDate.Month(), checkOutDate.Day(), 12, 0, 0, 0, time.Local)

	// Расчет стоимости
	costLabel := widget.NewLabel("Стоимость: -")
	advanceLabel := widget.NewLabel("")
	remainingLabel := widget.NewLabel("")

	updateCost := func() {
		if tariffSelect.SelectedIndex() >= 0 {
			// Use only the date part without time
			checkInDateOnly := time.Date(checkInDate.Year(), checkInDate.Month(), checkInDate.Day(), 0, 0, 0, 0, time.Local)
			checkOutDateOnly := time.Date(checkOutDate.Year(), checkOutDate.Month(), checkOutDate.Day(), 0, 0, 0, 0, time.Local)

			// Calculate days as the difference between dates
			days := int(checkOutDateOnly.Sub(checkInDateOnly).Hours()/24) + 1
			if days <= 0 {
				days = 1
			}

			// Calculate costs
			totalCost := float64(days) * tariffs[tariffSelect.SelectedIndex()].PricePerDay
			advancePayment := totalCost * 0.3 // 30% предоплата
			remainingAmount := totalCost - advancePayment

			// Create cost labels
			costLabel.SetText(fmt.Sprintf("Стоимость: %.2f руб. (%d дней)", totalCost, days))
			advanceLabel.SetText(fmt.Sprintf("Предоплата (30%%): %.2f руб.", advancePayment))
			remainingLabel.SetText(fmt.Sprintf("Остаток: %.2f руб.", remainingAmount))
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

	// Инициальный расчет стоимости
	updateCost()

	// Форма
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Домик", Widget: widget.NewLabel(cottage.Name)},
			{Text: "ФИО *", Widget: nameEntry},
			{Text: "Телефон *", Widget: phoneEntry},
			{Text: "Email", Widget: emailEntry},
			{Text: "Дата заезда", Widget: checkInPicker.button},
			{Text: "Дата выезда", Widget: checkOutPicker.button},
			{Text: "Тариф *", Widget: tariffSelect},
			{Text: "", Widget: costLabel},
			{Text: "", Widget: advanceLabel},
			{Text: "", Widget: remainingLabel},
			{Text: "Примечания", Widget: notesEntry},
		},
		OnSubmit: func() {
			// Валидация
			if nameEntry.Text == "" || phoneEntry.Text == "" || tariffSelect.SelectedIndex() < 0 {
				dialog.ShowError(fmt.Errorf("заполните обязательные поля (отмечены *)"), bc.window)
				return
			}

			if checkOutDate.Before(checkInDate) || checkOutDate.Equal(checkInDate) {
				dialog.ShowError(fmt.Errorf("дата выезда должна быть позже даты заезда"), bc.window)
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
			checkInDateOnly := time.Date(checkInDate.Year(), checkInDate.Month(), checkInDate.Day(), 0, 0, 0, 0, time.Local)
			checkOutDateOnly := time.Date(checkOutDate.Year(), checkOutDate.Month(), checkOutDate.Day(), 0, 0, 0, 0, time.Local)
			days := int(checkOutDateOnly.Sub(checkInDateOnly).Hours()/24) + 1
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

	d := dialog.NewCustom("Быстрое бронирование", "Отмена", form, bc.window)
	d.Resize(fyne.NewSize(500, 700))
	d.Show()
}

// createLegend создает легенду
func (bc *BookingCalendar) createLegend() fyne.CanvasObject {
	items := []fyne.CanvasObject{
		widget.NewLabel("Легенда:"),
		bc.createLegendItem("Свободно", GreenColor),
		bc.createLegendItem("Забронировано", color.NRGBA{R: 255, G: 193, B: 7, A: 255}), // Желтый
		bc.createLegendItem("Заселено", color.NRGBA{R: 255, G: 102, B: 0, A: 255}),      // Оранжевый
		widget.NewLabel("| Диагональ = Выезд/Заезд в один день"),
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
	case models.BookingStatusCancelled:
		return "Отменено"
	case models.BookingStatusTemporary:
		return "Временное"
	case models.BookingStatusBlocked:
		return "Заблокировано"
	case "completed":
		return "Завершено"
	default:
		return "Неизвестно"
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

// showEarlyCheckoutDialog показывает диалог раннего выселения
func (bc *BookingCalendar) showEarlyCheckoutDialog(booking *models.Booking, cottageName string) {
	// Создаем календарь для выбора новой даты выезда
	today := time.Now()
	checkInDate := time.Date(booking.CheckInDate.Year(), booking.CheckInDate.Month(), booking.CheckInDate.Day(), 0, 0, 0, 0, time.Local)
	originalCheckOut := time.Date(booking.CheckOutDate.Year(), booking.CheckOutDate.Month(), booking.CheckOutDate.Day(), 0, 0, 0, 0, time.Local)

	// Новая дата выезда должна быть между сегодняшним днем и оригинальной датой выезда
	var newCheckOutDate time.Time = today

	newCheckOutPicker := NewDatePickerButton(
		"Новая дата выезда",
		bc.window,
		func(t time.Time) {
			newCheckOutDate = t
		},
	)
	newCheckOutPicker.SetSelectedDate(today)

	// Расчет возврата денег
	refundLabel := widget.NewLabel("")

	updateRefund := func() {
		if !newCheckOutDate.IsZero() {
			// Рассчитываем сколько дней не использовано
			originalDays := int(originalCheckOut.Sub(checkInDate).Hours() / 24)
			usedDays := int(newCheckOutDate.Sub(checkInDate).Hours() / 24)

			if usedDays < 0 {
				usedDays = 0
			}
			if usedDays > originalDays {
				usedDays = originalDays
			}

			unusedDays := originalDays - usedDays
			refundAmount := (booking.TotalCost / float64(originalDays)) * float64(unusedDays)

			if refundAmount > 0 {
				refundLabel.SetText(fmt.Sprintf("К возврату: %.2f руб. (%d дней)", refundAmount, unusedDays))
			} else {
				refundLabel.SetText("Возврат не предусмотрен")
			}
		}
	}

	updateRefund()

	reasonEntry := widget.NewMultiLineEntry()
	reasonEntry.SetPlaceHolder("Причина раннего выезда...")
	reasonEntry.SetMinRowsVisible(3)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Гость", Widget: widget.NewLabel(booking.GuestName)},
			{Text: "Домик", Widget: widget.NewLabel(cottageName)},
			{Text: "Дата заезда", Widget: widget.NewLabel(booking.CheckInDate.Format("02.01.2006"))},
			{Text: "Планируемый выезд", Widget: widget.NewLabel(booking.CheckOutDate.Format("02.01.2006"))},
			{Text: "Новая дата выезда", Widget: newCheckOutPicker.button},
			{Text: "", Widget: refundLabel},
			{Text: "Причина", Widget: reasonEntry},
		},
		OnSubmit: func() {
			if newCheckOutDate.IsZero() {
				dialog.ShowError(fmt.Errorf("выберите дату выезда"), bc.window)
				return
			}

			if newCheckOutDate.Before(checkInDate) {
				dialog.ShowError(fmt.Errorf("дата выезда не может быть раньше даты заезда"), bc.window)
				return
			}

			// Устанавливаем время выезда на 12:00
			newCheckOutDateTime := time.Date(newCheckOutDate.Year(), newCheckOutDate.Month(), newCheckOutDate.Day(), 12, 0, 0, 0, time.Local)

			err := bc.bookingService.UpdateCheckOutDate(booking.ID, newCheckOutDateTime, reasonEntry.Text)
			if err != nil {
				dialog.ShowError(err, bc.window)
				return
			}

			bc.Update()
			dialog.ShowInformation("Успешно", "Дата выезда обновлена", bc.window)
		},
	}

	// Обновляем расчет при изменении даты
	newCheckOutPicker.button.OnTapped = func() {
		// Переопределяем обработчик после выбора даты
		originalTapped := newCheckOutPicker.button.OnTapped
		newCheckOutPicker.button.OnTapped = func() {
			originalTapped()
			updateRefund()
		}
	}

	d := dialog.NewCustom("Ранний выезд", "Отмена", form, bc.window)
	d.Resize(fyne.NewSize(450, 500))
	d.Show()
}
