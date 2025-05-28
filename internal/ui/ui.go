package ui

import (
	"fmt"
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
	w.Resize(fyne.NewSize(1200, 800))

	return &GuestApp{
		app:            a,
		window:         w,
		guestService:   guestService,
		cottageService: cottageService,
		tariffService:  tariffService,
		bookingService: bookingService,
	}
}

// Run запускает приложение
func (a *GuestApp) Run() {
	a.createUI()
	a.window.ShowAndRun()
}

// createUI формирует основное содержание окна с вкладками
func (a *GuestApp) createUI() {
	// Создаем виджеты календаря и списка бронирований
	calendarWidget := NewBookingCalendar(
		a.bookingService,
		a.cottageService,
		a.tariffService,
		a.window,
	)

	bookingListWidget := NewBookingListWidget(
		a.bookingService,
		a.cottageService,
		a.window,
	)

	// Устанавливаем взаимные обновления
	calendarWidget.SetOnRefresh(func() {
		bookingListWidget.loadData()
		if a.updateCottagesContent != nil {
			a.updateCottagesContent()
		}
	})

	bookingListWidget.SetOnRefresh(func() {
		calendarWidget.Update()
		if a.updateCottagesContent != nil {
			a.updateCottagesContent()
		}
	})

	// Создаем вкладки (убираем ненужные)
	tabs := container.NewAppTabs(
		container.NewTabItem("Календарь бронирований", calendarWidget),
		container.NewTabItem("Список бронирований", bookingListWidget),
		a.createTariffsTab(),
		a.createReportsTab(),
	)

	// Устанавливаем иконки для вкладок
	tabs.Items[0].Icon = theme.CalendarIcon()
	tabs.Items[1].Icon = theme.ListIcon()
	tabs.Items[2].Icon = theme.SettingsIcon()
	tabs.Items[3].Icon = theme.DocumentIcon()

	a.window.SetContent(tabs)
}

func (a *GuestApp) createTariffsTab() *container.TabItem {
	nameEntry := widget.NewEntry()
	priceEntry := widget.NewEntry()

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
			deleteBtn := hbox.Objects[2].(*widget.Button)

			nameLabel.SetText(tariff.Name)
			priceLabel.SetText(fmt.Sprintf("%.2f руб./день", tariff.PricePerDay))

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

	content := container.NewBorder(
		widget.NewCard("Добавить новый тариф", "", form),
		nil, nil, nil,
		widget.NewCard("Существующие тарифы", "", tariffList),
	)

	return container.NewTabItem("Тарифы", content)
}

// createReportsTab создает вкладку отчетов
func (a *GuestApp) createReportsTab() *container.TabItem {
	// Кнопки для генерации отчетов
	occupancyBtn := widget.NewButton("Отчет по заполняемости", func() {
		a.generateOccupancyReport()
	})

	revenueBtn := widget.NewButton("Отчет по доходам", func() {
		a.generateRevenueReport()
	})

	upcomingBtn := widget.NewButton("Предстоящие заезды", func() {
		a.showUpcomingArrivals()
	})

	content := container.NewVBox(
		widget.NewCard("Отчеты", "Выберите тип отчета для генерации",
			container.NewVBox(
				occupancyBtn,
				revenueBtn,
				upcomingBtn,
			),
		),
	)

	return container.NewTabItem("Отчеты",
		container.NewScroll(content),
	)
}

// generateOccupancyReport генерирует отчет по заполняемости
func (a *GuestApp) generateOccupancyReport() {
	// Выбор периода
	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now()

	bookings, err := a.bookingService.GetBookingsByDateRange(startDate, endDate)
	if err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	cottages, _ := a.cottageService.GetAllCottages()

	// Подсчет статистики
	totalDays := int(endDate.Sub(startDate).Hours() / 24)
	occupiedDays := make(map[int]int)

	for _, booking := range bookings {
		if booking.Status == models.BookingStatusCheckedIn ||
			booking.Status == models.BookingStatusCheckedOut {
			days := int(booking.CheckOutDate.Sub(booking.CheckInDate).Hours() / 24)
			occupiedDays[booking.CottageID] += days
		}
	}

	// Формируем отчет
	report := "ОТЧЕТ ПО ЗАПОЛНЯЕМОСТИ\n"
	report += fmt.Sprintf("Период: %s - %s\n\n",
		startDate.Format("02.01.2006"),
		endDate.Format("02.01.2006"))

	totalOccupancy := 0.0
	for _, cottage := range cottages {
		occupied := occupiedDays[cottage.ID]
		percentage := float64(occupied) / float64(totalDays) * 100
		totalOccupancy += percentage

		report += fmt.Sprintf("%s: %.1f%% (%d из %d дней)\n",
			cottage.Name, percentage, occupied, totalDays)
	}

	avgOccupancy := totalOccupancy / float64(len(cottages))
	report += fmt.Sprintf("\nСредняя заполняемость: %.1f%%\n", avgOccupancy)

	// Показываем отчет
	entry := widget.NewMultiLineEntry()
	entry.SetText(report)
	entry.Disable()

	d := dialog.NewCustom("Отчет по заполняемости", "Закрыть",
		container.NewScroll(entry), a.window)
	d.Resize(fyne.NewSize(500, 600))
	d.Show()
}

// generateRevenueReport генерирует отчет по доходам
func (a *GuestApp) generateRevenueReport() {
	// Аналогично генерируем отчет по доходам
	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now()

	bookings, err := a.bookingService.GetBookingsByDateRange(startDate, endDate)
	if err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	totalRevenue := 0.0
	revenueByMonth := make(map[string]float64)
	revenueByCottage := make(map[int]float64)

	for _, booking := range bookings {
		if booking.Status != models.BookingStatusCancelled {
			totalRevenue += booking.TotalCost

			monthKey := booking.CheckInDate.Format("01.2006")
			revenueByMonth[monthKey] += booking.TotalCost
			revenueByCottage[booking.CottageID] += booking.TotalCost
		}
	}

	report := "ОТЧЕТ ПО ДОХОДАМ\n"
	report += fmt.Sprintf("Период: %s - %s\n\n",
		startDate.Format("02.01.2006"),
		endDate.Format("02.01.2006"))

	report += fmt.Sprintf("Общий доход: %.2f руб.\n\n", totalRevenue)

	report += "По месяцам:\n"
	for month, revenue := range revenueByMonth {
		report += fmt.Sprintf("%s: %.2f руб.\n", month, revenue)
	}

	// Показываем отчет
	entry := widget.NewMultiLineEntry()
	entry.SetText(report)
	entry.Disable()

	d := dialog.NewCustom("Отчет по доходам", "Закрыть",
		container.NewScroll(entry), a.window)
	d.Resize(fyne.NewSize(500, 600))
	d.Show()
}

// showUpcomingArrivals показывает предстоящие заезды
func (a *GuestApp) showUpcomingArrivals() {
	bookings, err := a.bookingService.GetUpcomingBookings()
	if err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	cottages, _ := a.cottageService.GetAllCottages()
	cottageMap := make(map[int]string)
	for _, c := range cottages {
		cottageMap[c.ID] = c.Name
	}

	report := "ПРЕДСТОЯЩИЕ ЗАЕЗДЫ\n\n"

	for _, booking := range bookings {
		report += fmt.Sprintf("%s - %s:\n",
			booking.CheckInDate.Format("02.01.2006"),
			booking.CheckOutDate.Format("02.01.2006"))
		report += fmt.Sprintf("  Гость: %s\n", booking.GuestName)
		report += fmt.Sprintf("  Домик: %s\n", cottageMap[booking.CottageID])
		report += fmt.Sprintf("  Телефон: %s\n", booking.Phone)
		report += fmt.Sprintf("  Стоимость: %.2f руб.\n\n", booking.TotalCost)
	}

	if len(bookings) == 0 {
		report += "Нет предстоящих заездов\n"
	}

	// Показываем отчет
	entry := widget.NewMultiLineEntry()
	entry.SetText(report)
	entry.Disable()

	d := dialog.NewCustom("Предстоящие заезды", "Закрыть",
		container.NewScroll(entry), a.window)
	d.Resize(fyne.NewSize(500, 600))
	d.Show()
}
