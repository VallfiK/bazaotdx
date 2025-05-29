// internal/app/styled_app.go
package app

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/VallfIK/bazaotdx/internal/models"
	"github.com/VallfIK/bazaotdx/internal/service"
	"github.com/VallfIK/bazaotdx/internal/ui"
)

// StyledGuestApp — стилизованное приложение для учёта гостей
type StyledGuestApp struct {
	app                   fyne.App
	window                fyne.Window
	guestService          *service.GuestService
	cottageService        *service.CottageService
	tariffService         *service.TariffService
	bookingService        *service.BookingService
	updateCottagesContent func()
	cottages              []models.Cottage

	// UI виджеты
	calendarWidget    *ui.BookingCalendar
	bookingListWidget *ui.BookingListWidget
}

// NewStyledGuestApp создаёт новое стилизованное приложение
func NewStyledGuestApp(
	guestService *service.GuestService,
	cottageService *service.CottageService,
	tariffService *service.TariffService,
	bookingService *service.BookingService,
) *StyledGuestApp {
	a := app.New()

	// Устанавливаем кастомную тему
	a.Settings().SetTheme(&ui.ForestTheme{})

	w := a.NewWindow("🌲 Лесная База Отдыха - Система управления")
	w.Resize(fyne.NewSize(1600, 900))
	w.CenterOnScreen()

	app := &StyledGuestApp{
		app:            a,
		window:         w,
		guestService:   guestService,
		cottageService: cottageService,
		tariffService:  tariffService,
		bookingService: bookingService,
	}

	// Load cottages
	var err error
	app.cottages, err = app.cottageService.GetAllCottages()
	if err != nil {
		log.Printf("Error loading cottages: %v", err)
	}

	return app
}

// Run запускает приложение
func (a *StyledGuestApp) Run() {
	a.createStyledUI()
	a.window.ShowAndRun()
}

// createStyledUI формирует стилизованный интерфейс
func (a *StyledGuestApp) createStyledUI() {
	// Создаем верхнюю панель
	topBar := a.createTopBar()

	// Создаем календарный виджет
	a.calendarWidget = ui.NewBookingCalendar(
		a.bookingService,
		a.cottageService,
		a.tariffService,
		a.window,
	)

	// Создаем виджет списка бронирований
	a.bookingListWidget = ui.NewBookingListWidget(
		a.bookingService,
		a.cottageService,
		a.window,
	)

	// Устанавливаем взаимные обновления
	a.calendarWidget.SetOnRefresh(func() {
		a.bookingListWidget.Refresh()
	})

	a.bookingListWidget.SetOnRefresh(func() {
		a.calendarWidget.Update()
	})

	// Создаем боковую панель со статистикой
	sidePanel := a.createSidePanel()

	// Создаем основной контент
	mainContent := a.createMainContent()

	// Компонуем все вместе
	content := container.NewBorder(
		topBar,
		nil,
		sidePanel,
		nil,
		mainContent,
	)

	a.window.SetContent(content)
}

// createTopBar создает верхнюю панель
func (a *StyledGuestApp) createTopBar() fyne.CanvasObject {
	// Фон панели
	bg := canvas.NewRectangle(ui.DarkForestGreen)
	bg.SetMinSize(fyne.NewSize(0, 80))

	// Логотип и название
	title := canvas.NewText("🌲 Лесная База Отдыха", ui.Cream)
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Система управления бронированием", ui.LightCream)
	subtitle.TextSize = 14

	titleContainer := container.NewVBox(
		title,
		subtitle,
	)

	// Время и дата
	timeLabel := canvas.NewText("", ui.Cream)
	timeLabel.TextSize = 16
	updateTime := func() {
		now := time.Now()
		timeLabel.Text = now.Format("15:04:05\n02.01.2006")
		timeLabel.Refresh()
	}
	updateTime()
	go func() {
		for range time.Tick(1 * time.Second) {
			updateTime()
		}
	}()

	// Информация о пользователе
	userInfo := container.NewVBox(
		canvas.NewText("Администратор", ui.Cream),
		widget.NewButtonWithIcon("Выход", theme.LogoutIcon(), func() {
			dialog.ShowConfirm("Выход", "Вы уверены, что хотите выйти?",
				func(ok bool) {
					if ok {
						a.window.Close()
					}
				}, a.window)
		}),
	)

	// Компоновка
	content := container.NewBorder(
		nil, nil,
		container.NewPadded(titleContainer),
		container.NewHBox(
			container.NewPadded(timeLabel),
			widget.NewSeparator(),
			container.NewPadded(userInfo),
		),
		nil,
	)

	return container.NewMax(bg, content)
}

// createSidePanel создает боковую панель со статистикой
func (a *StyledGuestApp) createSidePanel() fyne.CanvasObject {
	// Статистические карточки
	statsCards := container.NewVBox()

	// Загружаем статистику
	updateStats := func() {
		statsCards.Objects = nil

		// Получаем статистику
		cottages, _ := a.cottageService.GetAllCottages()
		freeCottages := 0
		for _, c := range cottages {
			if c.Status == "free" {
				freeCottages++
			}
		}

		// Карточка свободных домиков
		freeCard := ui.InfoCard(
			theme.HomeIcon(),
			"Свободно домиков",
			fmt.Sprintf("%d из %d", freeCottages, len(cottages)),
			ui.ForestGreen,
		)
		freeCard.Resize(fyne.NewSize(250, 100))

		// Карточка текущего месяца
		now := time.Now()
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
		monthEnd := monthStart.AddDate(0, 1, -1)

		bookings, _ := a.bookingService.GetBookingsByDateRange(monthStart, monthEnd)
		activeBookings := 0
		totalRevenue := float64(0)

		for _, b := range bookings {
			if b.Status != models.BookingStatusCancelled {
				activeBookings++
				totalRevenue += b.TotalCost
			}
		}

		bookingCard := ui.InfoCard(
			theme.DocumentIcon(),
			"Бронирований в месяце",
			fmt.Sprintf("%d", activeBookings),
			ui.WoodBrown,
		)
		bookingCard.Resize(fyne.NewSize(250, 100))

		revenueCard := ui.InfoCard(
			theme.ContentCopyIcon(),
			"Доход за месяц",
			fmt.Sprintf("%.0f ₽", totalRevenue),
			ui.DarkForestGreen,
		)
		revenueCard.Resize(fyne.NewSize(250, 100))

		statsCards.Add(freeCard)
		statsCards.Add(bookingCard)
		statsCards.Add(revenueCard)
		statsCards.Refresh()
	}

	updateStats()

	// Кнопки быстрых действий
	quickActions := container.NewVBox(
		widget.NewCard("Быстрые действия", "", container.NewVBox(
			ui.NewGradientButton("➕ Новое бронирование", func() {
				a.showQuickBookingDialog()
			}),
			widget.NewButton("📊 Отчеты", func() {
				a.showReportsDialog()
			}),
			widget.NewButton("🏠 Управление домиками", func() {
				a.showCottagesDialog()
			}),
			widget.NewButton("💰 Тарифы", func() {
				a.showTariffsDialog()
			}),
		)),
	)

	// Обновление статистики каждые 30 секунд
	go func() {
		for range time.Tick(30 * time.Second) {
			updateStats()
		}
	}()

	panel := container.NewVBox(
		widget.NewCard("📊 Статистика", "", statsCards),
		widget.NewSeparator(),
		quickActions,
	)

	// Создаем скроллируемую панель
	scroll := container.NewVScroll(panel)
	scroll.SetMinSize(fyne.NewSize(280, 0))

	// Фон панели
	bg := canvas.NewRectangle(ui.Beige)

	return container.NewMax(bg, container.NewPadded(scroll))
}

// createMainContent создает основной контент
func (a *StyledGuestApp) createMainContent() fyne.CanvasObject {
	// Создаем вкладки с иконками и стилизованным видом
	calendarTab := container.NewTabItemWithIcon(
		"Календарь бронирований",
		theme.ViewFullScreenIcon(),
		container.NewHSplit(a.calendarWidget, a.bookingListWidget),
	)

	// Настраиваем разделитель
	if split, ok := calendarTab.Content.(*container.Split); ok {
		split.SetOffset(0.75)
	}

	tabs := container.NewAppTabs(calendarTab)

	return tabs
}

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

	notesEntry := widget.NewMultiLineEntry()
	notesEntry.SetPlaceHolder("Дополнительные примечания...")
	notesEntry.SetMinRowsVisible(3)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ФИО гостя", Widget: nameEntry},
			{Text: "Телефон", Widget: phoneEntry},
			{Text: "Email", Widget: emailEntry},
			{Text: "Домик", Widget: cottageSelect},
			{Text: "Тариф", Widget: tariffSelect},
			{Text: "Дата заезда", Widget: checkInPicker.button},
			{Text: "Дата выезда", Widget: checkOutPicker.button},
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
			price, err := strconv.ParseFloat(priceEntry.Text, 64)
			if err != nil {
				dialog.ShowError(fmt.Errorf("Неверный формат цены"), a.window)
				return
			}

			err = a.tariffService.UpdateTariff(tariff.ID, nameEntry.Text, price)
			if err != nil {
				dialog.ShowError(err, a.window)
				return
			}

			dialog.ShowInformation("Успешно", "Тариф обновлен", a.window)
		},
	}

	d := dialog.NewCustom("Редактировать тариф", "Отмена", form, a.window)
	d.Show()
}
