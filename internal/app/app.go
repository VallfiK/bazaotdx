// internal/app/app.go - Обновленная версия с применением темы
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

type GuestApp struct {
	fyneApp               fyne.App // Добавляем поле для Fyne приложения
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

func NewGuestApp(
	guestService *service.GuestService,
	cottageService *service.CottageService,
	tariffService *service.TariffService,
	bookingService *service.BookingService,
) *GuestApp {
	// Создаем Fyne приложение БЕЗ темы (тема будет установлена в main.go)
	a := app.New()
	w := a.NewWindow("🏕️ База отдыха - Система управления")
	w.Resize(fyne.NewSize(1400, 900))
	w.CenterOnScreen()

	app := &GuestApp{
		fyneApp:        a,
		window:         w,
		guestService:   guestService,
		cottageService: cottageService,
		tariffService:  tariffService,
		bookingService: bookingService,
	}

	// Загружаем домики
	var err error
	app.cottages, err = app.cottageService.GetAllCottages()
	if err != nil {
		log.Printf("Error loading cottages: %v", err)
	}

	return app
}

// SetFyneApp позволяет установить готовое Fyne приложение с темой
func (a *GuestApp) SetFyneApp(fyneApp fyne.App) {
	a.fyneApp = fyneApp
	// Пересоздаем окно с новым приложением
	a.window = fyneApp.NewWindow("🏕️ База отдыха - Система управления")
	a.window.Resize(fyne.NewSize(1400, 900))
	a.window.CenterOnScreen()
}

func (a *GuestApp) Run() {
	a.createStyledUI()
	a.window.ShowAndRun()
}

// createStyledUI создает интерфейс с применением новой темы
func (a *GuestApp) createStyledUI() {
	// Создаем статистические карточки
	statsContainer := a.createStatsContainer()

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

	// Создаем главный контейнер с разделенной панелью
	mainSplit := container.NewHSplit(a.calendarWidget, a.bookingListWidget)
	mainSplit.SetOffset(0.7) // 70% для календаря

	// Создаем вкладки с иконками и стилизацией
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Календарь", theme.ViewFullScreenIcon(),
			container.NewVBox(statsContainer, mainSplit)),
		container.NewTabItemWithIcon("Тарифы", theme.SettingsIcon(),
			a.createStyledTariffsTab()),
		container.NewTabItemWithIcon("Домики", theme.HomeIcon(),
			a.createStyledCottagesTab()),
		container.NewTabItemWithIcon("Отчеты", theme.DocumentIcon(),
			a.createStyledReportsTab()),
	)

	// Устанавливаем основное содержимое окна
	a.window.SetContent(tabs)
}

// createStatsContainer создает контейнер со статистикой
func (a *GuestApp) createStatsContainer() *fyne.Container {
	// Получаем статистику
	stats := a.calculateStats()

	// Создаем стилизованные карточки статистики
	totalCard := ui.CreateStatsCard("Всего домиков", fmt.Sprintf("%d", stats.TotalCottages), "в системе")
	occupiedCard := ui.CreateStatsCard("Занято", fmt.Sprintf("%d", stats.OccupiedCottages), "сегодня")
	upcomingCard := ui.CreateStatsCard("Предстоящие", fmt.Sprintf("%d", stats.UpcomingBookings), "заезды")
	revenueCard := ui.CreateStatsCard("Доход", fmt.Sprintf("₽ %.0f", stats.MonthlyRevenue), "за месяц")

	return container.NewGridWithColumns(4, totalCard, occupiedCard, upcomingCard, revenueCard)
}

// Stats структура для статистики
type Stats struct {
	TotalCottages    int
	OccupiedCottages int
	UpcomingBookings int
	MonthlyRevenue   float64
}

// calculateStats вычисляет статистику
func (a *GuestApp) calculateStats() Stats {
	stats := Stats{}

	// Общее количество домиков
	cottages, err := a.cottageService.GetAllCottages()
	if err == nil {
		stats.TotalCottages = len(cottages)
		for _, cottage := range cottages {
			if cottage.Status == "occupied" {
				stats.OccupiedCottages++
			}
		}
	}

	// Предстоящие бронирования
	upcomingBookings, err := a.bookingService.GetUpcomingBookings()
	if err == nil {
		stats.UpcomingBookings = len(upcomingBookings)
	}

	// Доходы за месяц
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	bookings, err := a.bookingService.GetBookingsByDateRange(startOfMonth, endOfMonth)
	if err == nil {
		for _, booking := range bookings {
			if booking.Status != models.BookingStatusCancelled {
				stats.MonthlyRevenue += booking.TotalCost
			}
		}
	}

	return stats
}

// createStyledTariffsTab создает стилизованную вкладку тарифов
func (a *GuestApp) createStyledTariffsTab() *fyne.Container {
	// Форма добавления тарифа
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Название тарифа")

	priceEntry := widget.NewEntry()
	priceEntry.SetPlaceHolder("Цена за день")

	// Стилизованная форма
	addForm := ui.CreateStyledCard("Добавить новый тариф", "",
		container.NewVBox(
			widget.NewForm(
				widget.NewFormItem("Название", nameEntry),
				widget.NewFormItem("Цена (руб./день)", priceEntry),
			),
			ui.NewPrimaryButton("Добавить тариф", func() {
				a.handleAddTariff(nameEntry.Text, priceEntry.Text, nameEntry, priceEntry)
			}),
		),
	)

	// Список тарифов
	tariffsList := a.createStyledTariffsList()

	return container.NewVBox(
		addForm,
		widget.NewSeparator(),
		tariffsList,
	)
}

// createStyledTariffsList создает стилизованный список тарифов
func (a *GuestApp) createStyledTariffsList() *fyne.Container {
	tariffs, err := a.tariffService.GetTariffs()
	if err != nil {
		return container.NewVBox(widget.NewLabel("Ошибка загрузки тарифов"))
	}

	tariffsContainer := container.NewVBox()

	for _, tariff := range tariffs {
		tariffCard := a.createTariffCard(tariff)
		tariffsContainer.Add(tariffCard)
	}

	return container.NewScroll(tariffsContainer)
}

// createTariffCard создает карточку тарифа
func (a *GuestApp) createTariffCard(tariff models.Tariff) *fyne.Container {
	// Создаем содержимое карточки
	content := container.NewHBox(
		widget.NewLabel(tariff.Name),
		widget.NewLabel(fmt.Sprintf("%.2f руб./день", tariff.PricePerDay)),
		ui.NewSecondaryButton("Изменить", func() {
			a.showEditTariffDialog(tariff)
		}),
		ui.NewSecondaryButton("Удалить", func() {
			a.showDeleteTariffConfirm(tariff)
		}),
	)

	return ui.CreateStyledCard("", "", content)
}

// createStyledCottagesTab создает стилизованную вкладку домиков
func (a *GuestApp) createStyledCottagesTab() *fyne.Container {
	// Форма добавления домика
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Название домика")

	addForm := ui.CreateStyledCard("Добавить новый домик", "",
		container.NewVBox(
			widget.NewForm(
				widget.NewFormItem("Название", nameEntry),
			),
			ui.NewPrimaryButton("Добавить домик", func() {
				a.handleAddCottage(nameEntry.Text, nameEntry)
			}),
		),
	)

	// Фильтры и поиск
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Поиск по названию...")

	statusFilter := widget.NewSelect([]string{"Все", "Свободные", "Занятые"}, nil)
	statusFilter.SetSelected("Все")

	filtersContainer := ui.CreateStyledCard("Фильтры", "",
		container.NewHBox(
			widget.NewLabel("Статус:"),
			statusFilter,
			widget.NewLabel("Поиск:"),
			searchEntry,
		),
	)

	// Список домиков
	cottagesList := a.createStyledCottagesList()

	return container.NewVBox(
		addForm,
		filtersContainer,
		widget.NewSeparator(),
		cottagesList,
	)
}

// createStyledCottagesList создает стилизованный список домиков
func (a *GuestApp) createStyledCottagesList() *fyne.Container {
	cottages, err := a.cottageService.GetAllCottages()
	if err != nil {
		return container.NewVBox(widget.NewLabel("Ошибка загрузки домиков"))
	}

	cottagesContainer := container.NewVBox()

	for _, cottage := range cottages {
		cottageCard := a.createCottageCard(cottage)
		cottagesContainer.Add(cottageCard)
	}

	return container.NewScroll(cottagesContainer)
}

// createCottageCard создает карточку домика
func (a *GuestApp) createCottageCard(cottage models.Cottage) *fyne.Container {
	// Определяем статус и цвет
	statusText := "Свободен"
	statusColor := ui.SuccessColor
	if cottage.Status == "occupied" {
		statusText = "Занят"
		statusColor = ui.WarmBrown
	}

	// Создаем индикатор статуса
	statusIndicator := ui.CreateStatusIndicator(cottage.Status)

	// Создаем содержимое карточки
	content := container.NewHBox(
		container.NewHBox(
			statusIndicator,
			widget.NewLabel(fmt.Sprintf("%d. %s", cottage.ID, cottage.Name)),
		),
		widget.NewLabel(statusText),
		ui.NewSecondaryButton("Изменить", func() {
			a.showEditCottageDialog(cottage)
		}),
		ui.NewSecondaryButton("Удалить", func() {
			a.showDeleteCottageConfirm(cottage)
		}),
	)

	return ui.CreateStyledCard("", "", content)
}

// createStyledReportsTab создает стилизованную вкладку отчетов
func (a *GuestApp) createStyledReportsTab() *fyne.Container {
	// Быстрые отчеты
	quickReports := ui.CreateStyledCard("Быстрые отчеты", "",
		container.NewVBox(
			ui.NewPrimaryButton("📊 Отчет о заполняемости", func() {
				a.generateOccupancyReport()
			}),
			ui.NewPrimaryButton("💰 Отчет о доходах", func() {
				a.generateRevenueReport()
			}),
			ui.NewPrimaryButton("👥 Предстоящие заезды", func() {
				a.showUpcomingArrivals()
			}),
		),
	)

	// Статистика
	statsCard := ui.CreateStyledCard("Статистика", "", a.createDetailedStats())

	return container.NewVBox(
		quickReports,
		statsCard,
	)
}

// createDetailedStats создает детальную статистику
func (a *GuestApp) createDetailedStats() *fyne.Container {
	stats := a.calculateStats()

	return container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Всего домиков: %d", stats.TotalCottages)),
		widget.NewLabel(fmt.Sprintf("Занято сегодня: %d", stats.OccupiedCottages)),
		widget.NewLabel(fmt.Sprintf("Предстоящие заезды: %d", stats.UpcomingBookings)),
		widget.NewLabel(fmt.Sprintf("Доход за месяц: %.2f руб.", stats.MonthlyRevenue)),
	)
}

// Обработчики событий
func (a *GuestApp) handleAddTariff(name, priceStr string, nameEntry, priceEntry *widget.Entry) {
	if name == "" || priceStr == "" {
		dialog.ShowError(fmt.Errorf("заполните все поля"), a.window)
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price <= 0 {
		dialog.ShowError(fmt.Errorf("неверный формат цены"), a.window)
		return
	}

	err = a.tariffService.CreateTariff(name, price)
	if err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	nameEntry.SetText("")
	priceEntry.SetText("")
	dialog.ShowInformation("Успешно", "Тариф создан", a.window)
}

func (a *GuestApp) handleAddCottage(name string, nameEntry *widget.Entry) {
	if name == "" {
		dialog.ShowError(fmt.Errorf("введите название домика"), a.window)
		return
	}

	err := a.cottageService.AddCottage(name)
	if err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	nameEntry.SetText("")
	if a.calendarWidget != nil {
		a.calendarWidget.Update()
	}
	dialog.ShowInformation("Успешно", "Домик добавлен", a.window)
}

func (a *GuestApp) showEditTariffDialog(tariff models.Tariff) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(tariff.Name)

	priceEntry := widget.NewEntry()
	priceEntry.SetText(fmt.Sprintf("%.2f", tariff.PricePerDay))

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Название", Widget: nameEntry},
			{Text: "Цена (руб./день)", Widget: priceEntry},
		},
		OnSubmit: func() {
			price, err := strconv.ParseFloat(priceEntry.Text, 64)
			if err != nil || price <= 0 {
				dialog.ShowError(fmt.Errorf("неверный формат цены"), a.window)
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

func (a *GuestApp) showDeleteTariffConfirm(tariff models.Tariff) {
	dialog.ShowConfirm("Удаление тарифа",
		fmt.Sprintf("Удалить тариф '%s'?", tariff.Name),
		func(confirmed bool) {
			if confirmed {
				err := a.tariffService.DeleteTariff(tariff.ID)
				if err != nil {
					dialog.ShowError(err, a.window)
					return
				}
				dialog.ShowInformation("Успешно", "Тариф удален", a.window)
			}
		}, a.window)
}

func (a *GuestApp) showEditCottageDialog(cottage models.Cottage) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(cottage.Name)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ID", Widget: widget.NewLabel(fmt.Sprintf("%d", cottage.ID))},
			{Text: "Название", Widget: nameEntry},
			{Text: "Статус", Widget: widget.NewLabel(cottage.Status)},
		},
		OnSubmit: func() {
			if nameEntry.Text == "" {
				dialog.ShowError(fmt.Errorf("название не может быть пустым"), a.window)
				return
			}

			err := a.cottageService.UpdateCottageName(cottage.ID, nameEntry.Text)
			if err != nil {
				dialog.ShowError(err, a.window)
				return
			}

			if a.calendarWidget != nil {
				a.calendarWidget.Update()
			}
			dialog.ShowInformation("Успешно", "Домик обновлен", a.window)
		},
	}

	d := dialog.NewCustom("Редактировать домик", "Отмена", form, a.window)
	d.Show()
}

func (a *GuestApp) showDeleteCottageConfirm(cottage models.Cottage) {
	if cottage.Status == "occupied" {
		dialog.ShowError(fmt.Errorf("невозможно удалить занятый домик"), a.window)
		return
	}

	dialog.ShowConfirm("Удаление домика",
		fmt.Sprintf("Удалить домик '%s'?", cottage.Name),
		func(confirmed bool) {
			if confirmed {
				err := a.cottageService.DeleteCottage(cottage.ID)
				if err != nil {
					dialog.ShowError(err, a.window)
					return
				}

				if a.calendarWidget != nil {
					a.calendarWidget.Update()
				}
				dialog.ShowInformation("Успешно", "Домик удален", a.window)
			}
		}, a.window)
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

		bookingCard := ui.CreateStyledCard(
			fmt.Sprintf("%s - %s", booking.GuestName, cottageName),
			fmt.Sprintf("Заезд: %s", booking.CheckInDate.Format("02.01.2006")),
			widget.NewLabel(fmt.Sprintf("Телефон: %s", booking.Phone)),
		)
		content.Add(bookingCard)
	}

	scroll := container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(500, 400))

	d := dialog.NewCustom("Предстоящие заезды", "Закрыть", scroll, a.window)
	d.Resize(fyne.NewSize(600, 500))
	d.Show()
}
