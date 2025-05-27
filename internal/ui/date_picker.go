package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// DatePickerButton — кастомная кнопка для выбора даты через календарь
type DatePickerButton struct {
	widget.BaseWidget
	selectedDate time.Time
	label        string
	button       *widget.Button
	onDateChange func(time.Time)
	window       fyne.Window
}

// NewDatePickerButton создает новую кнопку выбора даты
func NewDatePickerButton(label string, window fyne.Window, onChange func(time.Time)) *DatePickerButton {
	dpb := &DatePickerButton{
		label:        label,
		window:       window,
		onDateChange: onChange,
	}

	dpb.button = widget.NewButton(label, dpb.showCalendar)
	dpb.ExtendBaseWidget(dpb)
	return dpb
}

// showCalendar отображает календарь для выбора даты
func (dpb *DatePickerButton) showCalendar() {
	// Определяем начальную дату для календаря
	initialDate := time.Now()
	if !dpb.selectedDate.IsZero() {
		initialDate = dpb.selectedDate
	}

	// Создаем календарь
	cal := widget.NewCalendar(initialDate, func(t time.Time) {
		// Устанавливаем время в зависимости от типа даты
		var finalTime time.Time
		if dpb.label == "Дата заезда" {
			finalTime = time.Date(t.Year(), t.Month(), t.Day(), 14, 0, 0, 0, time.Local)
		} else {
			finalTime = time.Date(t.Year(), t.Month(), t.Day(), 12, 0, 0, 0, time.Local)
		}

		dpb.selectedDate = finalTime
		dpb.button.SetText(fmt.Sprintf("%s: %s", dpb.label, finalTime.Format("02.01.2006 15:04")))

		if dpb.onDateChange != nil {
			dpb.onDateChange(finalTime)
		}
	})

	// Создаем кнопки для управления
	todayBtn := widget.NewButton("Сегодня", func() {
		// Создаем новый календарь с сегодняшней датой
		today := time.Now()
		var finalTime time.Time
		if dpb.label == "Дата заезда" {
			finalTime = time.Date(today.Year(), today.Month(), today.Day(), 14, 0, 0, 0, time.Local)
		} else {
			finalTime = time.Date(today.Year(), today.Month(), today.Day(), 12, 0, 0, 0, time.Local)
		}

		dpb.selectedDate = finalTime
		dpb.button.SetText(fmt.Sprintf("%s: %s", dpb.label, finalTime.Format("02.01.2006 15:04")))

		if dpb.onDateChange != nil {
			dpb.onDateChange(finalTime)
		}
	})

	tomorrowBtn := widget.NewButton("Завтра", func() {
		// Создаем новый календарь с завтрашней датой
		tomorrow := time.Now().AddDate(0, 0, 1)
		var finalTime time.Time
		if dpb.label == "Дата заезда" {
			finalTime = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 14, 0, 0, 0, time.Local)
		} else {
			finalTime = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 12, 0, 0, 0, time.Local)
		}

		dpb.selectedDate = finalTime
		dpb.button.SetText(fmt.Sprintf("%s: %s", dpb.label, finalTime.Format("02.01.2006 15:04")))

		if dpb.onDateChange != nil {
			dpb.onDateChange(finalTime)
		}
	})

	clearBtn := widget.NewButton("Очистить", func() {
		dpb.selectedDate = time.Time{}
		dpb.button.SetText(dpb.label)
		if dpb.onDateChange != nil {
			dpb.onDateChange(time.Time{})
		}
	})

	// Создаем содержимое диалога
	content := container.NewVBox(
		cal,
		container.NewGridWithColumns(3, todayBtn, tomorrowBtn, clearBtn),
	)

	// Показываем диалог
	d := dialog.NewCustom("Выберите дату", "Закрыть", content, dpb.window)
	d.Resize(fyne.NewSize(300, 400))
	d.Show()
}

// GetSelectedDate возвращает выбранную дату
func (dpb *DatePickerButton) GetSelectedDate() time.Time {
	return dpb.selectedDate
}

// SetSelectedDate устанавливает дату программно
func (dpb *DatePickerButton) SetSelectedDate(t time.Time) {
	dpb.selectedDate = t
	if t.IsZero() {
		dpb.button.SetText(dpb.label)
	} else {
		dpb.button.SetText(fmt.Sprintf("%s: %s", dpb.label, t.Format("02.01.2006 15:04")))
	}
}

// CreateRenderer возвращает рендерер для виджета
func (dpb *DatePickerButton) CreateRenderer() fyne.WidgetRenderer {
	return dpb.button.CreateRenderer()
}
