// internal/app/styled_app.go
package app

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// showReportsDialog показывает диалог с отчетами
func (a *StyledGuestApp) showReportsDialog() {
	content := container.NewVBox(
		widget.NewCard("📊 Доступные отчеты", "", container.NewVBox(
			widget.NewButton("📈 Отчет о заполняемости", func() {
				dialog.ShowInformation("🔧 В разработке", "Эта функция будет доступна в следующей версии", a.window)
			}),
			widget.NewButton("💰 Финансовый отчет", func() {
				dialog.ShowInformation("🔧 В разработке", "Эта функция будет доступна в следующей версии", a.window)
			}),
			widget.NewButton("📋 Отчет по бронированиям", func() {
				dialog.ShowInformation("🔧 В разработке", "Эта функция будет доступна в следующей версии", a.window)
			}),
		)),
	)

	d := dialog.NewCustom("📊 Отчеты", "Закрыть", content, a.window)
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

			box.Objects[0].(*widget.Label).SetText(fmt.Sprintf("🏠 %d. %s", cottage.ID, cottage.Name))

			status := "🟢 Свободен"
			if cottage.Status == "occupied" {
				status = "🔴 Занят"
			}
			box.Objects[1].(*widget.Label).SetText(status)

			box.Objects[2].(*widget.Button).OnTapped = func() {
				a.showEditCottageDialog(cottage)
			}
		},
	)

	// Форма добавления с фиксированными размерами
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Название домика")
	nameEntry.Resize(fyne.NewSize(200, 40))

	addBtn := widget.NewButton("➕ Добавить домик", func() {
		if nameEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("Введите название"), a.window)
			return
		}

		err := a.cottageService.AddCottage(nameEntry.Text)
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}

		dialog.ShowInformation("✅ Успешно", "Домик добавлен", a.window)
		cottages, _ = a.cottageService.GetAllCottages()
		list.Refresh()
		a.calendarWidget.Update()
		nameEntry.SetText("")
	})

	content := container.NewBorder(
		widget.NewCard("➕ Добавить домик", "",
			container.NewBorder(nil, nil, nil, addBtn, nameEntry),
		),
		nil, nil, nil,
		list,
	)

	d := dialog.NewCustom("🏠 Управление домиками", "Закрыть", content, a.window)
	d.Resize(fyne.NewSize(600, 500))
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

			box.Objects[0].(*widget.Label).SetText(fmt.Sprintf("💰 %s", tariff.Name))
			box.Objects[1].(*widget.Label).SetText(fmt.Sprintf("%.0f ₽/сутки", tariff.PricePerDay))

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

	addBtn := widget.NewButton("➕ Добавить тариф", func() {
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

		dialog.ShowInformation("✅ Успешно", "Тариф добавлен", a.window)
		tariffs, _ = a.tariffService.GetTariffs()
		list.Refresh()
		nameEntry.SetText("")
		priceEntry.SetText("")
	})

	content := container.NewBorder(
		widget.NewCard("➕ Добавить тариф", "",
			container.NewVBox(
				nameEntry,
				priceEntry,
				addBtn,
			),
		),
		nil, nil, nil,
		list,
	)

	d := dialog.NewCustom("💰 Управление тарифами", "Закрыть", content, a.window)
	d.Resize(fyne.NewSize(500, 500))
	d.Show()
}
