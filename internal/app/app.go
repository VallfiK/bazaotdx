package app

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/VallfIK/bazaotdx/internal/models"
	"github.com/VallfIK/bazaotdx/internal/service"
	"github.com/VallfIK/bazaotdx/internal/ui"
)

// GuestApp — основное приложение для учёта гостей
type GuestApp struct {
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

// NewGuestApp создаёт новое приложение
func NewGuestApp(
	guestService *service.GuestService,
	cottageService *service.CottageService,
	tariffService *service.TariffService,
	bookingService *service.BookingService,
) *GuestApp {
	a := app.New()
	w := a.NewWindow("Учет гостей - База отдыха")
	w.Resize(fyne.NewSize(1400, 900))

	app := &GuestApp{
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
func (a *GuestApp) Run() {
	a.createUI()
	a.window.ShowAndRun()
}

// createUI формирует основное содержание окна с вкладками
func (a *GuestApp) createUI() {
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

	// Создаем разделенную панель для календаря и списка бронирований
	calendarTab := container.NewHSplit(a.calendarWidget, a.bookingListWidget)
	calendarTab.SetOffset(0.7) // 70% для календаря, 30% для списка

	// Создаем вкладки
	tabs := container.NewAppTabs(
		container.NewTabItem("Календарь", calendarTab),
		a.createTariffsTab(),
		a.createCottagesTab(),
		a.createReportsTab(),
	)

	// Устанавливаем основное содержимое окна
	a.window.SetContent(tabs)

	// Устанавливаем иконки для вкладок
	tabs.Items[0].Icon = theme.ViewFullScreenIcon()
	tabs.Items[1].Icon = theme.SettingsIcon()
	tabs.Items[2].Icon = theme.HomeIcon()
	tabs.Items[3].Icon = theme.DocumentIcon()
}

// createCottagesTab создает вкладку управления домиками
func (a *GuestApp) createCottagesTab() *container.TabItem {
	nameEntry := widget.NewEntry()
	nameEntry.PlaceHolder = "Введите название домика"

	// Список домиков
	var cottageList *widget.List
	var cottages []models.Cottage

	// Функция обновления списка домиков
	updateCottageList := func() {
		var err error
		cottages, err = a.cottageService.GetAllCottages()
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}
		a.cottages = cottages // Обновляем локальную копию
		if cottageList != nil {
			cottageList.Refresh()
		}
	}

	// Создаем список домиков
	cottageList = widget.NewList(
		func() int {
			return len(cottages)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Домик"),
				widget.NewLabel("Статус"),
				widget.NewButton("Удалить", nil),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id >= len(cottages) {
				return
			}
			cottage := cottages[id]

			hbox := item.(*fyne.Container)
			nameLabel := hbox.Objects[0].(*widget.Label)
			statusLabel := hbox.Objects[1].(*widget.Label)
			deleteBtn := hbox.Objects[2].(*widget.Button)

			nameLabel.SetText(cottage.Name)

			statusText := "Свободен"
			if cottage.Status == "occupied" {
				statusText = "Занят"
			}
			statusLabel.SetText(statusText)

			deleteBtn.OnTapped = func() {
				dialog.ShowConfirm("Подтверждение",
					fmt.Sprintf("Удалить домик '%s'?", cottage.Name),
					func(ok bool) {
						if ok {
							// Здесь нужно добавить метод удаления домика в CottageService
							// Пока что показываем предупреждение
							dialog.ShowInformation("Информация", "Функция удаления домиков будет добавлена позже", a.window)
						}
					}, a.window)
			}
		},
	)

	// Изначально загружаем домики
	updateCottageList()

	// Форма добавления нового домика
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Название домика", Widget: nameEntry},
		},
		OnSubmit: func() {
			if nameEntry.Text == "" {
				dialog.ShowError(fmt.Errorf("введите название домика"), a.window)
				return
			}

			err := a.cottageService.AddCottage(nameEntry.Text)
			if err != nil {
				dialog.ShowError(err, a.window)
				return
			}

			// Очищаем поле
			nameEntry.SetText("")

			updateCottageList()

			// Обновляем календарь после добавления домика
			if a.calendarWidget != nil {
				a.calendarWidget.Update()
			}

			dialog.ShowInformation("Успешно", "Домик добавлен", a.window)
		},
	}

	// Кнопка обновления
	refreshBtn := widget.NewButton("Обновить список", func() {
		updateCottageList()
		dialog.ShowInformation("Успешно", "Список обновлен", a.window)
	})

	content := container.NewVBox(
		form,
		refreshBtn,
		widget.NewLabel("Список домиков:"),
		cottageList,
	)

	return container.NewTabItem("Домики", content)
}

func (a *GuestApp) createTariffsTab() *container.TabItem {
	nameEntry := widget.NewEntry()
	nameEntry.PlaceHolder = "Введите название тарифа"

	priceEntry := widget.NewEntry()
	priceEntry.PlaceHolder = "Введите цену за день"

	// Список тарифов
	var tariffList *widget.List
	var tariffs []models.Tariff

	// Функция обновления списка тарифов
	updateTariffList := func() {
		var err error
		tariffs, err = a.tariffService.GetTariffs()
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}
		if tariffList != nil {
			tariffList.Refresh()
		}
	}

	// Создаем список тарифов
	tariffList = widget.NewList(
		func() int {
			return len(tariffs)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Тариф"),
				widget.NewLabel("Цена"),
				widget.NewButton("Изменить", nil),
				widget.NewButton("Удалить", nil),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id >= len(tariffs) {
				return
			}
			tariff := tariffs[id]

			hbox := item.(*fyne.Container)
			nameLabel := hbox.Objects[0].(*widget.Label)
			priceLabel := hbox.Objects[1].(*widget.Label)
			editBtn := hbox.Objects[2].(*widget.Button)
			deleteBtn := hbox.Objects[3].(*widget.Button)

			nameLabel.SetText(tariff.Name)
			priceLabel.SetText(fmt.Sprintf("%.2f руб./день", tariff.PricePerDay))

			editBtn.OnTapped = func() {
				a.showEditTariffDialog(tariff, updateTariffList)
			}

			deleteBtn.OnTapped = func() {
				dialog.ShowConfirm("Подтверждение",
					fmt.Sprintf("Удалить тариф '%s'?", tariff.Name),
					func(ok bool) {
						if ok {
							err := a.tariffService.DeleteTariff(tariff.ID)
							if err != nil {
								dialog.ShowError(err, a.window)
								return
							}
							updateTariffList()
							dialog.ShowInformation("Успешно", "Тариф удален", a.window)
						}
					}, a.window)
			}
		},
	)

	// Изначально загружаем тарифы
	updateTariffList()

	// Форма добавления нового тарифа
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Название тарифа", Widget: nameEntry},
			{Text: "Цена за день (руб.)", Widget: priceEntry},
		},
		OnSubmit: func() {
			if nameEntry.Text == "" || priceEntry.Text == "" {
				dialog.ShowError(fmt.Errorf("заполните все поля"), a.window)
				return
			}

			price, err := strconv.ParseFloat(priceEntry.Text, 64)
			if err != nil {
				dialog.ShowError(fmt.Errorf("неверный формат цены"), a.window)
				return
			}

			if price <= 0 {
				dialog.ShowError(fmt.Errorf("цена должна быть больше нуля"), a.window)
				return
			}

			err = a.tariffService.CreateTariff(nameEntry.Text, price)
			if err != nil {
				dialog.ShowError(err, a.window)
				return
			}

			// Очищаем поля
			nameEntry.SetText("")
			priceEntry.SetText("")

			updateTariffList()
			dialog.ShowInformation("Успешно", "Тариф создан", a.window)
		},
	}

	// Кнопка обновления
	refreshBtn := widget.NewButton("Обновить список", func() {
		updateTariffList()
		dialog.ShowInformation("Успешно", "Список обновлен", a.window)
	})

	content := container.NewVBox(
		form,
		refreshBtn,
		widget.NewLabel("Список тарифов:"),
		tariffList,
	)

	return container.NewTabItem("Тарифы", content)
}

// showEditTariffDialog показывает диалог редактирования тарифа
func (a *GuestApp) showEditTariffDialog(tariff models.Tariff, onUpdate func()) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(tariff.Name)

	priceEntry := widget.NewEntry()
	priceEntry.SetText(fmt.Sprintf("%.2f", tariff.PricePerDay))

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Название тарифа", Widget: nameEntry},
			{Text: "Цена за день (руб.)", Widget: priceEntry},
		},
		OnSubmit: func() {
			if nameEntry.Text == "" || priceEntry.Text == "" {
				dialog.ShowError(fmt.Errorf("заполните все поля"), a.window)
				return
			}

			price, err := strconv.ParseFloat(priceEntry.Text, 64)
			if err != nil {
				dialog.ShowError(fmt.Errorf("неверный формат цены"), a.window)
				return
			}

			if price <= 0 {
				dialog.ShowError(fmt.Errorf("цена должна быть больше нуля"), a.window)
				return
			}

			err = a.tariffService.UpdateTariff(tariff.ID, nameEntry.Text, price)
			if err != nil {
				dialog.ShowError(err, a.window)
				return
			}

			onUpdate()
			dialog.ShowInformation("Успешно", "Тариф обновлен", a.window)
		},
	}

	d := dialog.NewCustom("Редактировать тариф", "Отмена", form, a.window)
	d.Resize(fyne.NewSize(400, 300))
	d.Show()
}

func (a *GuestApp) createReportsTab() *container.TabItem {
	// Создаем кнопки для разных отчетов
	occupancyBtn := widget.NewButton("Отчет о заполняемости", func() {
		a.generateOccupancyReport()
	})

	revenueBtn := widget.NewButton("Отчет о доходах", func() {
		a.generateRevenueReport()
	})

	upcomingBtn := widget.NewButton("Предстоящие заезды", func() {
		a.showUpcomingArrivals()
	})

	// Статистика по текущему месяцу
	statsLabel := widget.NewLabel("Загрузка статистики...")
	a.updateStats(statsLabel)

	// Создаем контейнер с кнопками
	content := container.NewVBox(
		widget.NewCard("Статистика", "", statsLabel),
		widget.NewLabel("Отчеты:"),
		container.NewHBox(occupancyBtn, revenueBtn),
		upcomingBtn,
	)

	return container.NewTabItem("Отчеты", content)
}

// updateStats обновляет статистику
func (a *GuestApp) updateStats(label *widget.Label) {
	go func() {
		// Получаем статистику за текущий месяц
		now := time.Now()
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
		endOfMonth := startOfMonth.AddDate(0, 1, -1)

		bookings, err := a.bookingService.GetBookingsByDateRange(startOfMonth, endOfMonth)
		if err != nil {
			label.SetText("Ошибка загрузки статистики")
			return
		}

		totalBookings := len(bookings)
		var totalRevenue float64
		activeBookings := 0

		for _, booking := range bookings {
			if booking.Status != models.BookingStatusCancelled {
				totalRevenue += booking.TotalCost
				activeBookings++
			}
		}

		statsText := fmt.Sprintf(`Статистика за %s %d:
• Всего бронирований: %d
• Активных бронирований: %d  
• Общий доход: %.2f руб.
• Средний чек: %.2f руб.`,
			a.getMonthName(now.Month()), now.Year(),
			totalBookings, activeBookings, totalRevenue,
			func() float64 {
				if activeBookings > 0 {
					return totalRevenue / float64(activeBookings)
				}
				return 0
			}())

		label.SetText(statsText)
	}()
}

// getMonthName возвращает название месяца на русском
func (a *GuestApp) getMonthName(month time.Month) string {
	months := []string{
		"", "Январь", "Февраль", "Март", "Апрель", "Май", "Июнь",
		"Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь",
	}
	return months[month]
}

func (a *GuestApp) generateOccupancyReport() {
	dialog.ShowInformation("Информация", "Отчет о заполняемости будет реализован в следующей версии", a.window)
}

func (a *GuestApp) generateRevenueReport() {
	dialog.ShowInformation("Информация", "Отчет о доходах будет реализован в следующей версии", a.window)
}

func (a *GuestApp) showUpcomingArrivals() {
	bookings, err := a.bookingService.GetUpcomingBookings()
	if err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	if len(bookings) == 0 {
		dialog.ShowInformation("Информация", "Нет предстоящих заездов", a.window)
		return
	}

	content := container.NewVBox()
	for _, booking := range bookings {
		cottageName := ""
		for _, c := range a.cottages {
			if c.ID == booking.CottageID {
				cottageName = c.Name
				break
			}
		}

		card := widget.NewCard(
			fmt.Sprintf("%s - %s", booking.GuestName, cottageName),
			fmt.Sprintf("Заезд: %s", booking.CheckInDate.Format("02.01.2006")),
			widget.NewLabel(fmt.Sprintf("Телефон: %s", booking.Phone)),
		)
		content.Add(card)
	}

	scroll := container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(500, 400))

	d := dialog.NewCustom("Предстоящие заезды", "Закрыть", scroll, a.window)
	d.Resize(fyne.NewSize(600, 500))
	d.Show()
}
