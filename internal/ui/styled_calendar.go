// internal/ui/styled_calendar.go
package ui

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// StyledCalendarCell создает стилизованную ячейку календаря
type StyledCalendarCell struct {
	widget.BaseWidget
	date      time.Time
	status    string
	guestName string
	onTapped  func()
	isToday   bool
	isWeekend bool
	cottageID int
}

func NewStyledCalendarCell(date time.Time, status, guestName string, onTapped func()) *StyledCalendarCell {
	cell := &StyledCalendarCell{
		date:      date,
		status:    status,
		guestName: guestName,
		onTapped:  onTapped,
		isToday:   isToday(date),
		isWeekend: isWeekend(date),
	}
	cell.ExtendBaseWidget(cell)
	return cell
}

func (c *StyledCalendarCell) CreateRenderer() fyne.WidgetRenderer {
	// Определяем цвет фона в зависимости от статуса
	var bgColor color.Color
	switch c.status {
	case "free":
		if c.isWeekend {
			bgColor = MintGreen
		} else {
			bgColor = LightCream
		}
	case "booked":
		bgColor = color.NRGBA{R: 255, G: 193, B: 7, A: 180} // Янтарный
	case "occupied":
		bgColor = ForestGreen
	default:
		bgColor = LightCream
	}

	// Создаем фон с закругленными углами
	bg := canvas.NewRectangle(bgColor)
	bg.CornerRadius = 6

	// Рамка для сегодняшнего дня
	var border *canvas.Rectangle
	if c.isToday {
		border = canvas.NewRectangle(color.Transparent)
		border.StrokeColor = DarkForestGreen
		border.StrokeWidth = 3
		border.CornerRadius = 6
	}

	// Текст даты
	dateText := canvas.NewText(fmt.Sprintf("%d", c.date.Day()), DarkBrown)
	dateText.TextSize = 14
	dateText.TextStyle = fyne.TextStyle{Bold: true}
	dateText.Alignment = fyne.TextAlignCenter

	// Текст гостя
	var guestText *canvas.Text
	if c.guestName != "" {
		displayName := c.guestName
		if len(displayName) > 8 {
			displayName = displayName[:8] + "..."
		}
		guestText = canvas.NewText(displayName, color.White)
		guestText.TextSize = 10
		guestText.Alignment = fyne.TextAlignCenter
	}

	// Иконка статуса
	var statusIcon *canvas.Image
	switch c.status {
	case "free":
		if c.isWeekend {
			// Можно добавить специальную иконку для выходных
		}
	case "booked":
		statusIcon = canvas.NewImageFromResource(theme.ConfirmIcon())
		statusIcon.FillMode = canvas.ImageFillContain
	case "occupied":
		statusIcon = canvas.NewImageFromResource(theme.HomeIcon())
		statusIcon.FillMode = canvas.ImageFillContain
	}

	return &styledCalendarCellRenderer{
		bg:         bg,
		border:     border,
		dateText:   dateText,
		guestText:  guestText,
		statusIcon: statusIcon,
		cell:       c,
	}
}

type styledCalendarCellRenderer struct {
	bg         *canvas.Rectangle
	border     *canvas.Rectangle
	dateText   *canvas.Text
	guestText  *canvas.Text
	statusIcon *canvas.Image
	cell       *StyledCalendarCell
}

func (r *styledCalendarCellRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)

	if r.border != nil {
		r.border.Resize(size)
	}

	// Размещаем дату в верхнем углу
	r.dateText.Move(fyne.NewPos(5, 5))
	r.dateText.Resize(fyne.NewSize(size.Width-10, 20))

	// Размещаем имя гостя в центре
	if r.guestText != nil {
		r.guestText.Move(fyne.NewPos(5, size.Height/2-10))
		r.guestText.Resize(fyne.NewSize(size.Width-10, 20))
	}

	// Размещаем иконку в правом нижнем углу
	if r.statusIcon != nil {
		iconSize := float32(16)
		r.statusIcon.Move(fyne.NewPos(size.Width-iconSize-5, size.Height-iconSize-5))
		r.statusIcon.Resize(fyne.NewSize(iconSize, iconSize))
	}
}

func (r *styledCalendarCellRenderer) MinSize() fyne.Size {
	return fyne.NewSize(80, 60)
}

func (r *styledCalendarCellRenderer) Refresh() {
	// Обновляем цвета и тексты
	r.dateText.Text = fmt.Sprintf("%d", r.cell.date.Day())
	r.dateText.Refresh()

	if r.guestText != nil && r.cell.guestName != "" {
		displayName := r.cell.guestName
		if len(displayName) > 8 {
			displayName = displayName[:8] + "..."
		}
		r.guestText.Text = displayName
		r.guestText.Refresh()
	}
}

func (r *styledCalendarCellRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.bg}

	if r.border != nil {
		objects = append(objects, r.border)
	}

	objects = append(objects, r.dateText)

	if r.guestText != nil {
		objects = append(objects, r.guestText)
	}

	if r.statusIcon != nil {
		objects = append(objects, r.statusIcon)
	}

	return objects
}

func (r *styledCalendarCellRenderer) Destroy() {}

func (c *StyledCalendarCell) Tapped(_ *fyne.PointEvent) {
	if c.onTapped != nil {
		c.onTapped()
	}
}

func (c *StyledCalendarCell) TappedSecondary(_ *fyne.PointEvent) {}

// StyledMonthHeader создает стилизованный заголовок месяца
func StyledMonthHeader(month time.Time, onPrev, onNext, onToday func()) fyne.CanvasObject {
	// Фон заголовка
	bg := canvas.NewRectangle(ForestGreen)
	bg.CornerRadius = 10

	// Название месяца
	monthName := getMonthNameRu(month.Month())
	yearText := fmt.Sprintf("%s %d", monthName, month.Year())

	monthLabel := canvas.NewText(yearText, Cream)
	monthLabel.TextSize = 20
	monthLabel.TextStyle = fyne.TextStyle{Bold: true}
	monthLabel.Alignment = fyne.TextAlignCenter

	// Кнопки навигации
	prevBtn := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), onPrev)
	nextBtn := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), onNext)
	todayBtn := widget.NewButton("Сегодня", onToday)

	// Компоновка
	content := container.NewBorder(
		nil, nil,
		prevBtn,
		container.NewHBox(todayBtn, nextBtn),
		container.NewCenter(monthLabel),
	)

	return container.NewMax(bg, container.NewPadded(content))
}

// StyledDayHeader создает стилизованный заголовок дня недели
func StyledDayHeader(weekday string) fyne.CanvasObject {
	bg := canvas.NewRectangle(LightBrown)
	bg.CornerRadius = 4

	text := canvas.NewText(weekday, Cream)
	text.TextSize = 12
	text.TextStyle = fyne.TextStyle{Bold: true}
	text.Alignment = fyne.TextAlignCenter

	return container.NewMax(bg, container.NewCenter(text))
}

// StyledLegend создает стилизованную легенду
func StyledLegend() fyne.CanvasObject {
	items := []fyne.CanvasObject{
		createLegendItem("Свободно", LightCream),
		createLegendItem("Забронировано", color.NRGBA{R: 255, G: 193, B: 7, A: 180}),
		createLegendItem("Заселено", ForestGreen),
		createLegendItem("Выходной", MintGreen),
	}

	legend := container.NewHBox(items...)

	// Обрамляем в карточку
	card := widget.NewCard("", "", legend)

	return card
}

func createLegendItem(text string, color color.Color) fyne.CanvasObject {
	// Цветной квадрат
	colorBox := canvas.NewRectangle(color)
	colorBox.CornerRadius = 4
	colorBox.SetMinSize(fyne.NewSize(20, 20))

	// Текст
	label := widget.NewLabel(text)

	return container.NewHBox(colorBox, label)
}

// Вспомогательные функции
func isToday(date time.Time) bool {
	now := time.Now()
	return date.Year() == now.Year() && date.YearDay() == now.YearDay()
}

func isWeekend(date time.Time) bool {
	weekday := date.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

func getMonthNameRu(month time.Month) string {
	months := []string{
		"", "Январь", "Февраль", "Март", "Апрель", "Май", "Июнь",
		"Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь",
	}
	return months[month]
}
