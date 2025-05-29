// internal/ui/styled_calendar.go
package ui

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/VallfIK/bazaotdx/internal/models"
)

// StyledCalendarCell - стилизованная ячейка календаря
type StyledCalendarCell struct {
	widget.BaseWidget

	status    string
	guestName string
	cottageID int
	date      time.Time
	onTapped  func()

	// Визуальные элементы
	background *canvas.Rectangle
	text       *canvas.Text
	indicator  *canvas.Circle
}

func NewStyledCalendarCell(status, guestName string, cottageID int, date time.Time, onTapped func()) *StyledCalendarCell {
	cell := &StyledCalendarCell{
		status:    status,
		guestName: guestName,
		cottageID: cottageID,
		date:      date,
		onTapped:  onTapped,
	}

	cell.ExtendBaseWidget(cell)
	cell.createVisualElements()

	return cell
}

func (c *StyledCalendarCell) createVisualElements() {
	// Создаем фон
	c.background = canvas.NewRectangle(GetStatusColor(c.status))
	c.background.CornerRadius = 6
	c.background.StrokeWidth = 1
	c.background.StrokeColor = BorderLight

	// Создаем текст
	displayText := "+"
	if c.guestName != "" {
		if len(c.guestName) > 8 {
			displayText = c.guestName[:6] + ".."
		} else {
			displayText = c.guestName
		}
	}

	c.text = canvas.NewText(displayText, TextDark)
	c.text.Alignment = fyne.TextAlignCenter
	c.text.TextSize = 10

	// Создаем индикатор статуса
	c.indicator = canvas.NewCircle(GetStatusColor(c.status))
	c.indicator.Resize(fyne.NewSize(6, 6))
}

func (c *StyledCalendarCell) Tapped(_ *fyne.PointEvent) {
	if c.onTapped != nil {
		c.onTapped()
	}
}

func (c *StyledCalendarCell) TappedSecondary(_ *fyne.PointEvent) {}

func (c *StyledCalendarCell) CreateRenderer() fyne.WidgetRenderer {
	return &styledCalendarCellRenderer{cell: c}
}

func (c *StyledCalendarCell) MinSize() fyne.Size {
	return fyne.NewSize(35, 60)
}

// Рендерер для стилизованной ячейки
type styledCalendarCellRenderer struct {
	cell *StyledCalendarCell
}

func (r *styledCalendarCellRenderer) Layout(size fyne.Size) {
	r.cell.background.Resize(size)

	// Позиционируем текст по центру
	textSize := r.cell.text.MinSize()
	r.cell.text.Move(fyne.NewPos(
		(size.Width-textSize.Width)/2,
		(size.Height-textSize.Height)/2,
	))

	// Позиционируем индикатор в правом верхнем углу
	r.cell.indicator.Move(fyne.NewPos(size.Width-8, 2))
}

func (r *styledCalendarCellRenderer) MinSize() fyne.Size {
	return r.cell.MinSize()
}

func (r *styledCalendarCellRenderer) Refresh() {
	r.cell.background.FillColor = GetStatusColor(r.cell.status)
	r.cell.indicator.FillColor = GetStatusColor(r.cell.status)

	canvas.Refresh(r.cell.background)
	canvas.Refresh(r.cell.text)
	canvas.Refresh(r.cell.indicator)
}

func (r *styledCalendarCellRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		r.cell.background,
		r.cell.text,
		r.cell.indicator,
	}
}

func (r *styledCalendarCellRenderer) Destroy() {}

// StyledBookingCard - стилизованная карточка бронирования
type StyledBookingCard struct {
	widget.BaseWidget

	booking  models.Booking
	cottage  models.Cottage
	onTapped func()

	// Визуальные элементы
	background      *canvas.Rectangle
	nameText        *canvas.Text
	dateText        *canvas.Text
	phoneText       *canvas.Text
	statusIndicator *canvas.Circle
}

func NewStyledBookingCard(booking models.Booking, cottage models.Cottage, onTapped func()) *StyledBookingCard {
	card := &StyledBookingCard{
		booking:  booking,
		cottage:  cottage,
		onTapped: onTapped,
	}

	card.ExtendBaseWidget(card)
	card.createVisualElements()

	return card
}

func (c *StyledBookingCard) createVisualElements() {
	// Фон карточки
	c.background = canvas.NewRectangle(color.White)
	c.background.CornerRadius = 8
	c.background.StrokeWidth = 1

	// Определяем цвет границы в зависимости от статуса
	switch c.booking.Status {
	case models.BookingStatusBooked:
		c.background.StrokeColor = SuccessColor
	case models.BookingStatusCheckedIn:
		c.background.StrokeColor = WarmBrown
	default:
		c.background.StrokeColor = BorderLight
	}

	// Создаем текстовые элементы
	c.nameText = canvas.NewText(c.booking.GuestName, TextDark)
	c.nameText.TextStyle = fyne.TextStyle{Bold: true}
	c.nameText.TextSize = 12

	dateText := fmt.Sprintf("%s - %s • %s",
		c.booking.CheckInDate.Format("02.01"),
		c.booking.CheckOutDate.Format("02.01"),
		c.cottage.Name,
	)
	c.dateText = canvas.NewText(dateText, TextLight)
	c.dateText.TextSize = 10

	c.phoneText = canvas.NewText(c.booking.Phone, TextLight)
	c.phoneText.TextSize = 9

	// Индикатор статуса
	statusColor := SuccessColor
	if c.booking.Status == models.BookingStatusCheckedIn {
		statusColor = WarmBrown
	}
	c.statusIndicator = canvas.NewCircle(statusColor)
}

func (c *StyledBookingCard) Tapped(_ *fyne.PointEvent) {
	if c.onTapped != nil {
		c.onTapped()
	}
}

func (c *StyledBookingCard) TappedSecondary(_ *fyne.PointEvent) {}

func (c *StyledBookingCard) CreateRenderer() fyne.WidgetRenderer {
	return &styledBookingCardRenderer{card: c}
}

func (c *StyledBookingCard) MinSize() fyne.Size {
	return fyne.NewSize(250, 80)
}

// Рендерер для карточки бронирования
type styledBookingCardRenderer struct {
	card *StyledBookingCard
}

func (r *styledBookingCardRenderer) Layout(size fyne.Size) {
	r.card.background.Resize(size)

	// Позиционируем элементы
	padding := float32(8)

	// Индикатор статуса
	r.card.statusIndicator.Resize(fyne.NewSize(8, 8))
	r.card.statusIndicator.Move(fyne.NewPos(padding, padding))

	// Текстовые элементы
	textX := padding + 16
	r.card.nameText.Move(fyne.NewPos(textX, padding))
	r.card.dateText.Move(fyne.NewPos(textX, padding+20))
	r.card.phoneText.Move(fyne.NewPos(textX, padding+40))
}

func (r *styledBookingCardRenderer) MinSize() fyne.Size {
	return r.card.MinSize()
}

func (r *styledBookingCardRenderer) Refresh() {
	canvas.Refresh(r.card.background)
	canvas.Refresh(r.card.nameText)
	canvas.Refresh(r.card.dateText)
	canvas.Refresh(r.card.phoneText)
	canvas.Refresh(r.card.statusIndicator)
}

func (r *styledBookingCardRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		r.card.background,
		r.card.statusIndicator,
		r.card.nameText,
		r.card.dateText,
		r.card.phoneText,
	}
}

func (r *styledBookingCardRenderer) Destroy() {}

// Стилизованная форма бронирования
func CreateStyledBookingForm(window fyne.Window, onSubmit func(models.Booking)) *fyne.Container {
	// Заголовок формы
	title := canvas.NewText("Новое бронирование", PrimaryGreen)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.TextSize = 18
	title.Alignment = fyne.TextAlignCenter

	// Поля формы
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("ФИО гостя")

	phoneEntry := widget.NewEntry()
	phoneEntry.SetPlaceHolder("+7 (999) 123-45-67")

	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("email@example.com")

	// Создаем стилизованные labels
	nameLabel := canvas.NewText("ФИО гостя *", TextDark)
	nameLabel.TextSize = 12

	phoneLabel := canvas.NewText("Телефон *", TextDark)
	phoneLabel.TextSize = 12

	emailLabel := canvas.NewText("Email", TextDark)
	emailLabel.TextSize = 12

	// Группируем поля
	nameGroup := container.NewVBox(nameLabel, nameEntry)
	phoneGroup := container.NewVBox(phoneLabel, phoneEntry)
	emailGroup := container.NewVBox(emailLabel, emailEntry)

	// Кнопки
	submitBtn := NewPrimaryButton("Создать бронирование", func() {
		// Валидация и создание бронирования
		if nameEntry.Text == "" || phoneEntry.Text == "" {
			// Показать ошибку
			return
		}

		booking := models.Booking{
			GuestName: nameEntry.Text,
			Phone:     phoneEntry.Text,
			Email:     emailEntry.Text,
		}

		if onSubmit != nil {
			onSubmit(booking)
		}
	})

	cancelBtn := NewSecondaryButton("Отмена", func() {
		window.Close()
	})

	// Создаем фон формы
	formBg := canvas.NewRectangle(color.White)
	formBg.CornerRadius = 12
	formBg.StrokeWidth = 1
	formBg.StrokeColor = BorderLight

	// Собираем форму
	formContent := container.NewVBox(
		title,
		widget.NewSeparator(),
		nameGroup,
		phoneGroup,
		emailGroup,
		widget.NewSeparator(),
		container.NewHBox(cancelBtn, submitBtn),
	)

	// Добавляем отступы
	paddedForm := container.NewPadded(formContent)

	return container.NewStack(formBg, paddedForm)
}

// Функция для создания стилизованной статистической карточки
func CreateStatsCard(title, value, subtitle string) *fyne.Container {
	// Создаем фон
	bg := canvas.NewRectangle(color.White)
	bg.CornerRadius = 12
	bg.StrokeWidth = 1
	bg.StrokeColor = BorderLight

	// Создаем текстовые элементы
	valueText := canvas.NewText(value, SecondaryGreen)
	valueText.TextStyle = fyne.TextStyle{Bold: true}
	valueText.TextSize = 24
	valueText.Alignment = fyne.TextAlignCenter

	titleText := canvas.NewText(title, TextDark)
	titleText.TextSize = 12
	titleText.Alignment = fyne.TextAlignCenter

	subtitleText := canvas.NewText(subtitle, TextLight)
	subtitleText.TextSize = 10
	subtitleText.Alignment = fyne.TextAlignCenter

	// Создаем content
	content := container.NewVBox(
		container.NewCenter(valueText),
		container.NewCenter(titleText),
		container.NewCenter(subtitleText),
	)

	return container.NewStack(bg, container.NewPadded(content))
}
