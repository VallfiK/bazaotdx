package app

import (
	"fmt"
	"image/color"
	"log"
	"strconv"
	"strings"
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
	calendarWidget        *ui.BookingCalendar
	bookingListWidget     *ui.BookingListWidget
	topBar                fyne.CanvasObject
	sidePanel             fyne.CanvasObject
	mainContent           fyne.CanvasObject
}

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

	// Initialize widgets
	app.calendarWidget = ui.NewBookingCalendar(
		app.bookingService,
		app.cottageService,
		app.tariffService,
		app.window,
	)
	app.bookingListWidget = ui.NewBookingListWidget(
		app.bookingService,
		app.cottageService,
		app.window,
	)

	// Load cottages
	var err error
	app.cottages, err = app.cottageService.GetAllCottages()
	if err != nil {
		log.Printf("Error loading cottages: %v", err)
	}

	app.createStyledUI()

	return app
}

// Run запускает приложение
func (a *StyledGuestApp) Run() {
	a.window.ShowAndRun()
}

// createStyledUI формирует стилизованный интерфейс
func (a *StyledGuestApp) createStyledUI() {
	// Создаем верхнюю панель
	topBar := a.createTopBar()

	// Создаем календарный виджет
	calendarWidget := ui.NewBookingCalendar(
		a.bookingService,
		a.cottageService,
		a.tariffService,
		a.window,
	)

	// Создаем виджет списка бронирований
	bookingListWidget := ui.NewBookingListWidget(
		a.bookingService,
		a.cottageService,
		a.window,
	)

	sidePanel := a.createSidePanel()
	mainContent := a.createMainContent()

	// Устанавливаем минимальный размер окна
	a.window.Resize(fyne.NewSize(1200, 800))
	// Запрещаем изменение размера окна
	a.window.SetFixedSize(true)

	// Создаем основное окно с контейнером
	content := container.NewBorder(
		topBar, nil,
		sidePanel,
		mainContent,
	)

	// Сохраняем ссылки на виджеты
	a.topBar = topBar
	a.calendarWidget = calendarWidget
	a.bookingListWidget = bookingListWidget
	a.sidePanel = sidePanel
	a.mainContent = mainContent

	a.window.SetContent(content)
}

// createTopBar создает верхнюю панель
func (a *StyledGuestApp) createTopBar() fyne.CanvasObject {
	// Фон панели
	bg := canvas.NewRectangle(ui.DarkForestGreen)
	bg.SetMinSize(fyne.NewSize(0, 80))

	// Логотип и название
	title := canvas.NewText("🌲 Лесная База Отдыха", color.White)
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Система управления бронированием", color.White)
	subtitle.TextSize = 14

	titleContainer := container.NewVBox(
		title,
		subtitle,
	)

	// Время и дата
	timeLabel := canvas.NewText("", color.White)
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
	adminText := canvas.NewText("Администратор", color.White)
logoutBtn := widget.NewButtonWithIcon("Выход", theme.LogoutIcon(), func() {
		dialog.ShowConfirm("Выход", "Вы уверены, что хотите выйти?",
			func(ok bool) {
				if ok {
					a.window.Close()
				}
			}, a.window)
	})
userInfo := container.NewVBox(adminText, logoutBtn)

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

		// Создаем карточки
		statsCards.Add(a.createStatsDisplay())
	}

	// Запускаем первое обновление
	updateStats()

	// Создаем кнопки быстрого доступа
	quickActions := container.NewVBox(
		widget.NewButtonWithIcon("Быстрое бронирование", theme.ContentAddIcon(), func() {
			a.showQuickBookingDialog()
		}),
		widget.NewButtonWithIcon("Предстоящие заезды", theme.CalendarIcon(), func() {
			a.showUpcomingArrivals()
		}),
	)

	return container.NewVBox(
		statsCards,
		widget.NewSeparator(),
		quickActions,
	)
}

// createMainContent создает основной контент
func (a *StyledGuestApp) createMainContent() fyne.CanvasObject {
	// Создаем разделенную панель для календаря и списка бронирований
	calendarTab := container.NewHSplit(a.calendarWidget, a.bookingListWidget)
	calendarTab.SetOffset(0.7) // 70% для календаря, 30% для списка

	// Создаем вкладки
	tabs := container.NewAppTabs(
		container.NewTabItem("Календарь", calendarTab),
		a.createTariffsTab(),
		a.createCottagesTab(),
	)

	// Устанавливаем иконки для вкладок
	tabs.Items[0].Icon = theme.ViewFullScreenIcon()
	tabs.Items[1].Icon = theme.SettingsIcon()
	tabs.Items[2].Icon = theme.HomeIcon()

	// Устанавливаем фиксированный размер для вкладок
	tabs.Resize(fyne.NewSize(1000, 600))
	return tabs
}

// createStatsDisplay создает панель со статистикой
func (a *StyledGuestApp) createStatsDisplay() fyne.CanvasObject {
	// Создаем виджеты для статистики
	statsLabel := widget.NewLabel("")
	updateBtn := widget.NewButton("Обновить", func() {
		a.updateStats(statsLabel)
	})

	// Изначально обновляем статистику
	go a.updateStats(statsLabel)

	// Создаем контейнер
	content := container.NewVBox(
		widget.NewCard("Статистика", "", statsLabel),
		updateBtn,
	)

	return content
}

// createCottagesTab создает вкладку домиков
func (a *StyledGuestApp) createCottagesTab() *container.TabItem {
	// Создаем список домиков
	var cottages []models.Cottage
	var cottageList *widget.List

	// Функция обновления списка домиков
	var updateCottageList func()
	updateCottageList = func() {
		c, err := a.cottageService.GetAllCottages()
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}
		cottages = c

		if cottageList != nil {
			cottageList.Refresh()
		}
	}

	// Инициализируем список
	cottageList = widget.NewList(
		func() int {
			return len(cottages)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Название"),
				widget.NewLabel("Статус"),
				widget.NewButton("Изменить", nil),
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
			editBtn := hbox.Objects[2].(*widget.Button)

			nameLabel.SetText(cottage.Name)
			statusLabel.SetText(cottage.Status)

			editBtn.OnTapped = func() {
				a.showEditCottageDialog(cottage)
			}
		},
	)

	// Форма добавления нового домика
	nameEntry := widget.NewEntry()
	statusSelect := widget.NewSelect([]string{"active", "inactive"}, nil)

	addForm := widget.NewForm(
		widget.NewFormItem("Название", nameEntry),
		widget.NewFormItem("Статус", statusSelect),
	)

	addBtn := widget.NewButton("Добавить домик", func() {
		if nameEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("введите название"), a.window)
			return
		}

		err := a.cottageService.AddCottage(nameEntry.Text)
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}

		// Очищаем поля ввода
		nameEntry.SetText("")
		statusSelect.Selected = "active"
		updateCottageList()

		dialog.ShowInformation("Успешно",
			fmt.Sprintf("Домик '%s' добавлен", nameEntry.Text),
			a.window,
		)
	})

	// Кнопка обновления
	refreshBtn := widget.NewButton("Обновить список", func() {
		updateCottageList()
		dialog.ShowInformation("Успешно", "Список обновлен", a.window)
	})

	// Изначально загружаем домики
	updateCottageList()

	// Создаем контейнер со скроллом для списка
	listScroll := container.NewScroll(cottageList)
	listScroll.SetMinSize(fyne.NewSize(600, 400))

	content := container.NewBorder(
		container.NewVBox(
			widget.NewCard("Добавить новый домик", "", container.NewVBox(addForm, addBtn)),
			container.NewHBox(refreshBtn),
			widget.NewSeparator(),
		),
		nil, nil, nil,
		listScroll,
	)

	return container.NewTabItem("Домики", content)
}

// createTariffsTab создает вкладку тарифов
func (a *StyledGuestApp) createTariffsTab() *container.TabItem {
	// Форма поиска
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Поиск по названию...")

	// Список тарифов
	var filteredTariffs []models.Tariff
	var tariffList *widget.List

	// Создаем виджеты
	statsLabel := widget.NewLabel("")

	// Функция обновления списка тарифов
	var updateTariffList func()
	updateTariffList = func() {
		tariffs, err := a.tariffService.GetTariffs()
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}

		// Фильтруем тарифы по названию
		filteredTariffs = make([]models.Tariff, 0)
		for _, tariff := range tariffs {
			if searchEntry.Text != "" && !strings.Contains(tariff.Name, searchEntry.Text) {
				continue
			}
			filteredTariffs = append(filteredTariffs, tariff)
		}

		if tariffList != nil {
			tariffList.Refresh()
		}

		// Обновляем статистику
		total := len(tariffs)
		statsLabel.SetText(fmt.Sprintf("Всего тарифов: %d", total))
	}

	// Инициализируем список
	tariffList = widget.NewList(
		func() int {
			return len(filteredTariffs)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Название"),
				widget.NewLabel("Цена"),
				widget.NewButton("Изменить", nil),
				widget.NewButton("Удалить", nil),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id >= len(filteredTariffs) {
				return
			}
			tariff := filteredTariffs[id]

			hbox := item.(*fyne.Container)
			nameLabel := hbox.Objects[0].(*widget.Label)
			priceLabel := hbox.Objects[1].(*widget.Label)
			editBtn := hbox.Objects[2].(*widget.Button)
			deleteBtn := hbox.Objects[3].(*widget.Button)

			nameLabel.SetText(tariff.Name)
			priceLabel.SetText(fmt.Sprintf("%.2f руб./день", tariff.PricePerDay))

			editBtn.OnTapped = func() {
				a.showEditTariffDialog(tariff)
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
					},
					a.window)
			}
		},
	)

	// Форма поиска
	searchEntry.OnChanged = func(text string) {
		updateTariffList()
	}

	// Форма добавления нового тарифа
	nameEntry := widget.NewEntry()
	priceEntry := widget.NewEntry()

	addForm := widget.NewForm(
		widget.NewFormItem("Название", nameEntry),
		widget.NewFormItem("Цена за сутки", priceEntry),
	)

	addBtn := widget.NewButton("Добавить тариф", func() {
		if nameEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("введите название"), a.window)
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

		// Очищаем поля ввода
		nameEntry.SetText("")
		priceEntry.SetText("")
		updateTariffList()

		dialog.ShowInformation("Успешно",
			fmt.Sprintf("Тариф '%s' добавлен", nameEntry.Text),
			a.window,
		)
	})

	// Кнопка обновления
	refreshBtn := widget.NewButton("Обновить список", func() {
		updateTariffList()
		dialog.ShowInformation("Успешно", "Список обновлен", a.window)
	})

	// Изначально загружаем тарифы
	updateTariffList()

	// Создаем контейнер со скроллом для списка
	listScroll := container.NewScroll(tariffList)
	listScroll.SetMinSize(fyne.NewSize(600, 400))

	content := container.NewBorder(
		container.NewVBox(
			searchEntry,
			widget.NewCard("Добавить новый тариф", "", container.NewVBox(addForm, addBtn)),
			container.NewHBox(refreshBtn, statsLabel),
			widget.NewSeparator(),
		),
		nil, nil, nil,
		listScroll,
	)

	// Возвращаем контейнер с вкладкой
	return container.NewTabItem("Тарифы", content)
}



func (a *StyledGuestApp) createReportsTab() *container.TabItem {
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
func (a *StyledGuestApp) updateStats(label *widget.Label) {
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
func (a *StyledGuestApp) getMonthName(month time.Month) string {
	months := []string{
		"", "Январь", "Февраль", "Март", "Апрель", "Май", "Июнь",
		"Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь",
	}
	return months[month]
}

func (a *StyledGuestApp) generateOccupancyReport() {
	dialog.ShowInformation("Информация", "Отчет о заполняемости будет реализован в следующей версии", a.window)
}

func (a *StyledGuestApp) generateRevenueReport() {
	dialog.ShowInformation("Информация", "Отчет о доходах будет реализован в следующей версии", a.window)
}

func (a *StyledGuestApp) showUpcomingArrivals() {
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
