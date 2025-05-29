// internal/ui/styled_app.go
package app

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/VallfIK/bazaotdx/internal/models"
	"github.com/VallfIK/bazaotdx/internal/ui"
)









// showQuickBookingDialog показывает диалог быстрого бронирования
func (a *StyledGuestApp) showQuickBookingDialog() {
	// Получаем список свободных домиков
	cottages, err := a.cottageService.GetFreeCottages()
	if err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	if len(cottages) == 0 {
		dialog.ShowInformation("Информация", "Нет свободных домиков", a.window)
		return
	}

	// Получаем тарифы
	tariffs, err := a.tariffService.GetTariffs()
	if err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	// Создаем форму
	nameEntry := ui.StyledEntry("ФИО гостя")
	phoneEntry := ui.StyledEntry("+7 (999) 123-45-67")
	emailEntry := ui.StyledEntry("email@example.com")

	cottageOptions := make([]string, len(cottages))
	for i, c := range cottages {
		cottageOptions[i] = fmt.Sprintf("%d. %s", c.ID, c.Name)
	}
	cottageSelect := widget.NewSelect(cottageOptions, nil)

	tariffOptions := make([]string, len(tariffs))
	for i, t := range tariffs {
		tariffOptions[i] = fmt.Sprintf("%s - %.2f ₽/сутки", t.Name, t.PricePerDay)
	}
	tariffSelect := widget.NewSelect(tariffOptions, nil)

	// Даты
	checkInDate := time.Now().Add(24 * time.Hour)
	checkOutDate := checkInDate.Add(24 * time.Hour)

	checkInPicker := ui.NewDatePickerButton("Дата заезда", a.window, func(t time.Time) {
			checkInDate = t
		})
	checkInPicker.SetSelectedDate(checkInDate)
	checkOutPicker := ui.NewDatePickerButton("Дата выезда", a.window, func(t time.Time) {
			checkOutDate = t
		})
	checkOutPicker.SetSelectedDate(checkOutDate)

	notesEntry := ui.StyledEntry("Дополнительные примечания...")
	notesEntry.SetPlaceHolder("Дополнительные примечания...")
	notesEntry.SetMinRowsVisible(3)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ФИО гостя", Widget: nameEntry},
			{Text: "Телефон", Widget: phoneEntry},
			{Text: "Email", Widget: emailEntry},
			{Text: "Домик", Widget: cottageSelect},
			{Text: "Тариф", Widget: tariffSelect},
			{Text: "Дата заезда", Widget: checkInPicker},
			{Text: "Дата выезда", Widget: checkOutPicker},
			{Text: "Примечания", Widget: notesEntry},
		},
		OnSubmit: func() {
			// Валидация
			if nameEntry.Text == "" || phoneEntry.Text == "" ||
				cottageSelect.SelectedIndex() < 0 || tariffSelect.SelectedIndex() < 0 {
				dialog.ShowError(fmt.Errorf("Заполните все обязательные поля"), a.window)
				return
			}

			// Создаем бронирование
			booking := models.Booking{
				CottageID:    cottages[cottageSelect.SelectedIndex()].ID,
				GuestName:    nameEntry.Text,
				Phone:        phoneEntry.Text,
				Email:        emailEntry.Text,
				CheckInDate:  checkInPicker.GetSelectedDate(),
				CheckOutDate: checkOutPicker.GetSelectedDate(),
				TariffID:     tariffs[tariffSelect.SelectedIndex()].ID,
				Notes:        notesEntry.Text,
			}

			_, err := a.bookingService.CreateBooking(booking)
			if err != nil {
				dialog.ShowError(err, a.window)
				return
			}

			a.calendarWidget.Update()
			dialog.ShowInformation("Успешно", "Бронирование создано", a.window)
		},
		OnCancel: func() {},
	}

	d := dialog.NewCustom("Новое бронирование", "Отмена", form, a.window)
	d.Resize(fyne.NewSize(500, 600))
	d.Show()
}

// showReportsDialog показывает диалог с отчетами
func (a *StyledGuestApp) showReportsDialog() {
	content := container.NewVBox(
		widget.NewCard("Доступные отчеты", "", container.NewVBox(
			widget.NewButton("📊 Отчет о заполняемости", func() {
				dialog.ShowInformation("В разработке", "Эта функция будет доступна в следующей версии", a.window)
			}),
			widget.NewButton("💰 Финансовый отчет", func() {
				dialog.ShowInformation("В разработке", "Эта функция будет доступна в следующей версии", a.window)
			}),
			widget.NewButton("📅 Отчет по бронированиям", func() {
				dialog.ShowInformation("В разработке", "Эта функция будет доступна в следующей версии", a.window)
			}),
		)),
	)

	d := dialog.NewCustom("Отчеты", "Закрыть", content, a.window)
	d.Resize(fyne.NewSize(400, 300))
	d.Show()
}

// showCottagesDialog показывает диалог управления домиками
func (a *StyledGuestApp) showCottagesDialog() {
	// Создаем список домиков
	cottages, _ := a.cottageService.GetAllCottages()

	list := widget.NewList(
		func() int { return len(cottages) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Домик"),
				widget.NewLabel("Статус"),
				widget.NewButton("Изменить", nil),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			cottage := cottages[id]
			box := item.(*fyne.Container)

			box.Objects[0].(*widget.Label).SetText(fmt.Sprintf("%d. %s", cottage.ID, cottage.Name))

			status := "Свободен"
			if cottage.Status == "occupied" {
				status = "Занят"
			}
			box.Objects[1].(*widget.Label).SetText(status)

			box.Objects[2].(*widget.Button).OnTapped = func() {
				a.showEditCottageDialog(cottage)
			}
		},
	)

	// Форма добавления
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Название домика")

	addBtn := widget.NewButton("Добавить домик", func() {
		if nameEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("Введите название"), a.window)
			return
		}

		err := a.cottageService.AddCottage(nameEntry.Text)
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}

		dialog.ShowInformation("Успешно", "Домик добавлен", a.window)
		cottages, _ = a.cottageService.GetAllCottages()
		list.Refresh()
		a.calendarWidget.Update()
	})

	content := container.NewBorder(
		widget.NewCard("Добавить домик", "",
			container.NewBorder(nil, nil, nil, addBtn, nameEntry),
		),
		nil, nil, nil,
		list,
	)

	d := dialog.NewCustom("Управление домиками", "Закрыть", content, a.window)
	d.Resize(fyne.NewSize(600, 500))
	d.Show()
}

// showEditCottageDialog показывает диалог редактирования домика
func (a *StyledGuestApp) showEditCottageDialog(cottage models.Cottage) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(cottage.Name)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ID", Widget: widget.NewLabel(fmt.Sprintf("%d", cottage.ID))},
			{Text: "Название", Widget: nameEntry},
			{Text: "Статус", Widget: widget.NewLabel(cottage.Status)},
		},
		OnSubmit: func() {
			err := a.cottageService.UpdateCottageName(cottage.ID, nameEntry.Text)
			if err != nil {
				dialog.ShowError(err, a.window)
				return
			}
			dialog.ShowInformation("Успешно", "Домик обновлен", a.window)
			a.calendarWidget.Update()
		},
	}

	d := dialog.NewCustom("Редактировать домик", "Отмена", form, a.window)
	d.Show()
}

// showTariffsDialog показывает диалог управления тарифами
func (a *StyledGuestApp) showTariffsDialog() {
	tariffs, _ := a.tariffService.GetTariffs()

	list := widget.NewList(
		func() int { return len(tariffs) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Тариф"),
				widget.NewLabel("Цена"),
				widget.NewButton("Изменить", nil),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			tariff := tariffs[id]
			box := item.(*fyne.Container)

			box.Objects[0].(*widget.Label).SetText(tariff.Name)
			box.Objects[1].(*widget.Label).SetText(fmt.Sprintf("%.2f ₽/сутки", tariff.PricePerDay))

			box.Objects[2].(*widget.Button).OnTapped = func() {
				a.showEditTariffDialog(tariff)
			}
		},
	)

	// Форма добавления
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Название тарифа")

	priceEntry := widget.NewEntry()
	priceEntry.SetPlaceHolder("Цена за сутки")

	addBtn := widget.NewButton("Добавить тариф", func() {
		if nameEntry.Text == "" || priceEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("Заполните все поля"), a.window)
			return
		}

		price, err := strconv.ParseFloat(priceEntry.Text, 64)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Неверный формат цены"), a.window)
			return
		}

		err = a.tariffService.CreateTariff(nameEntry.Text, price)
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}

		dialog.ShowInformation("Успешно", "Тариф добавлен", a.window)
		tariffs, _ = a.tariffService.GetTariffs()
		list.Refresh()
	})

	content := container.NewBorder(
		widget.NewCard("Добавить тариф", "",
			container.NewVBox(
				nameEntry,
				priceEntry,
				addBtn,
			),
		),
		nil, nil, nil,
		list,
	)

	d := dialog.NewCustom("Управление тарифами", "Закрыть", content, a.window)
	d.Resize(fyne.NewSize(500, 500))
	d.Show()
}

// showEditTariffDialog показывает диалог редактирования тарифа
func (a *StyledGuestApp) showEditTariffDialog(tariff models.Tariff) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(tariff.Name)

	priceEntry := widget.NewEntry()
	priceEntry.SetText(fmt.Sprintf("%.2f", tariff.PricePerDay))

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Название", Widget: nameEntry},
			{Text: "Цена за сутки", Widget: priceEntry},
		},
		OnSubmit: func() {
			if nameEntry.Text == "" {
				dialog.ShowError(fmt.Errorf("введите название"), a.window)
				return
			}

			price, err := strconv.ParseFloat(priceEntry.Text, 64)
			if err != nil {
				dialog.ShowError(fmt.Errorf("неверный формат цены"), a.window)
				return
			}

			err = a.tariffService.UpdateTariff(tariff.ID, nameEntry.Text, price)
			if err != nil {
				dialog.ShowError(err, a.window)
				return
			}

			dialog.ShowInformation("Успешно", "Тариф сохранен", a.window)
		},
	}

	d := dialog.NewCustom("Редактирование тарифа", "Отмена", form, a.window)
	d.Resize(fyne.NewSize(500, 400))
	d.Show()
}
