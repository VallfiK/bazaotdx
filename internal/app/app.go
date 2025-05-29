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

// StyledGuestApp — стилизованное приложение для учёта гостей "Звуки Леса"
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

	w := a.NewWindow("🌲 Звуки Леса - База Отдыха")

	// ФИКСИРУЕМ РАЗМЕР ОКНА
	fixedSize := fyne.NewSize(1400, 800)
	w.Resize(fixedSize)
	w.SetFixedSize(true) // Запрещаем изменение размера
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
	// Создаем верхнюю панель с фиксированной высотой
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

	// Создаем боковую панель с фиксированной шириной
	sidePanel := a.createSidePanel()
	mainContent := a.createMainContent()

	// Сохраняем ссылки на виджеты
	a.topBar = topBar
	a.calendarWidget = calendarWidget
	a.bookingListWidget = bookingListWidget
	a.sidePanel = sidePanel
	a.mainContent = mainContent

	// Создаем основное окно с контейнером и фиксированными размерами
	content := container.NewBorder(
		topBar,                                // Верх - фиксированная высота 80px
		nil,                                   // Низ
		container.NewWithoutLayout(sidePanel), // Левая панель - фиксированная ширина 280px
		nil,                                   // Правая панель
		mainContent,                           // Центр - оставшееся место
	)

	a.window.SetContent(content)
}

// createTopBar создает верхнюю панель с фиксированной высотой
func (a *StyledGuestApp) createTopBar() fyne.CanvasObject {
	// Фон панели с градиентом
	bg := canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		ratio := float32(x) / float32(w)
		// Градиент от темно-зеленого к лесному зеленому
		return lerpColor(ui.DarkForestGreen, ui.ForestGreen, ratio)
	})
	bg.SetMinSize(fyne.NewSize(0, 80))

	// Логотип и название с улучшенной типографикой
	title := canvas.NewText("🌲 Звуки Леса", color.White)
	title.TextSize = 32
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("База отдыха • Система управления", ui.GoldenYellow)
	subtitle.TextSize = 16
	subtitle.TextStyle = fyne.TextStyle{Italic: true}

	titleContainer := container.NewVBox(
		title,
		subtitle,
	)

	// Время и дата с улучшенным стилем
	timeLabel := canvas.NewText("", color.White)
	timeLabel.TextSize = 18
	timeLabel.TextStyle = fyne.TextStyle{Bold: true}
	timeLabel.Alignment = fyne.TextAlignCenter

	updateTime := func() {
		now := time.Now()
		fyne.Do(func() {
			timeLabel.Text = now.Format("15:04:05\n02.01.2006")
			timeLabel.Refresh()
		})
	}

	updateTime()
	go func() {
		for range time.Tick(1 * time.Second) {
			fyne.Do(func() {
				updateTime()
			})
		}
	}()

	// Информация о пользователе
	adminText := canvas.NewText("👤 Администратор", color.White)
	adminText.TextSize = 14
	adminText.TextStyle = fyne.TextStyle{Bold: true}

	logoutBtn := widget.NewButtonWithIcon("Выход", theme.LogoutIcon(), func() {
		dialog.ShowConfirm("Выход из системы", "Вы уверены, что хотите выйти из системы управления?",
			func(ok bool) {
				if ok {
					a.window.Close()
				}
			}, a.window)
	})
	logoutBtn.Importance = widget.DangerImportance

	userInfo := container.NewVBox(adminText, logoutBtn)

	// Компоновка с фиксированными размерами
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

// createSidePanel создает боковую панель со статистикой с фиксированной шириной
func (a *StyledGuestApp) createSidePanel() fyne.CanvasObject {
	// Контейнер с фиксированной шириной
	sideContainer := container.NewWithoutLayout()

	// Статистические карточки
	statsContainer := container.NewVBox()

	// Функция обновления статистики с ФИКСИРОВАННЫМИ размерами
	updateStats := func() {
		statsContainer.Objects = nil
		statsContainer.Add(a.createFixedStatsDisplay())
		statsContainer.Refresh()
	}

	// Запускаем первое обновление
	updateStats()

	// Создаем кнопки быстрого доступа с фиксированными размерами
	quickBookingBtn := widget.NewButtonWithIcon("✨ Быстрое бронирование", theme.ContentAddIcon(), func() {
		a.showQuickBookingDialogFromSidePanel()
	})
	quickBookingBtn.Importance = widget.HighImportance
	quickBookingBtn.Resize(fyne.NewSize(260, 40))

	upcomingBtn := widget.NewButtonWithIcon("📅 Предстоящие заезды", theme.MailAttachmentIcon(), func() {
		a.showUpcomingArrivals()
	})
	upcomingBtn.Resize(fyne.NewSize(260, 40))

	quickActions := container.NewVBox(
		quickBookingBtn,
		upcomingBtn,
	)

	// Основной контент боковой панели
	mainSideContent := container.NewVBox(
		statsContainer,
		widget.NewCard("", "", widget.NewSeparator()),
		widget.NewCard("Быстрые действия", "", quickActions),
	)

	// Размещаем в контейнере с фиксированной позицией
	sideContainer.Add(mainSideContent)
	mainSideContent.Move(fyne.NewPos(0, 0))
	mainSideContent.Resize(fyne.NewSize(280, 700))

	return sideContainer
}

// createFixedStatsDisplay создает панель со статистикой с ФИКСИРОВАННЫМИ размерами
func (a *StyledGuestApp) createFixedStatsDisplay() fyne.CanvasObject {
	// Создаем виджеты для статистики с фиксированными размерами
	statsLabel := widget.NewRichTextFromMarkdown("Загрузка статистики...")
	statsLabel.Resize(fyne.NewSize(260, 120)) // Увеличиваем ширину для предотвращения переноса текста
	statsLabel.Wrapping = fyne.TextWrapWord   // Добавляем перенос текста по словам

	// Кнопка обновления с фиксированным размером
	updateBtn := widget.NewButtonWithIcon("🔄 Обновить", theme.ViewRefreshIcon(), func() {
		// Блокируем изменение размеров во время обновления
		a.updateStats(statsLabel)
	})
	updateBtn.Resize(fyne.NewSize(260, 40)) // ФИКСИРОВАННЫЙ РАЗМЕР
	updateBtn.Importance = widget.MediumImportance

	// Изначально обновляем статистику
	go a.updateStats(statsLabel)

	// Создаем контейнер с фиксированными размерами
	content := container.NewVBox(
		widget.NewCard("📊 Статистика", "", statsLabel),
		updateBtn,
	)

	// ПРИНУДИТЕЛЬНО устанавливаем размер контейнера
	content.Resize(fyne.NewSize(270, 200))

	return content
}

// createMainContent создает основной контент с фиксированными размерами
func (a *StyledGuestApp) createMainContent() fyne.CanvasObject {
	// Создаем разделенную панель для календаря и списка бронирований
	calendarTab := container.NewHSplit(a.calendarWidget, a.bookingListWidget)
	calendarTab.SetOffset(0.7) // 70% для календаря, 30% для списка

	// Создаем вкладки
	tabs := container.NewAppTabs(
		container.NewTabItem("📅 Календарь", calendarTab),
		a.createTariffsTab(),
		a.createCottagesTab(),
	)

	// ФИКСИРУЕМ размер для вкладок
	tabs.Resize(fyne.NewSize(1100, 700))

	return tabs
}

// createCottagesTab создает вкладку домиков с фиксированными размерами
func (a *StyledGuestApp) createCottagesTab() *container.TabItem {
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

	// Инициализируем список с фиксированными размерами
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
				a.showEditCottageDialogFixed(cottage)
			}
		},
	)

	// Форма добавления нового домика
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Название домика")

	addForm := widget.NewForm(
		widget.NewFormItem("🏠 Название", nameEntry),
	)

	addBtn := widget.NewButtonWithIcon("Добавить домик", theme.ContentAddIcon(), func() {
		if nameEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("введите название"), a.window)
			return
		}

		err := a.cottageService.AddCottage(nameEntry.Text)
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}

		nameEntry.SetText("")
		updateCottageList()

		dialog.ShowInformation("Успешно",
			fmt.Sprintf("Домик '%s' добавлен", nameEntry.Text),
			a.window,
		)
	})
	addBtn.Importance = widget.HighImportance

	// Кнопка обновления с фиксированным размером
	refreshBtn := widget.NewButtonWithIcon("🔄 Обновить", theme.ViewRefreshIcon(), func() {
		updateCottageList()
		dialog.ShowInformation("Успешно", "Список обновлен", a.window)
	})
	refreshBtn.Resize(fyne.NewSize(150, 40))

	// Изначально загружаем домики
	updateCottageList()

	// Создаем контейнер со скроллом для списка с фиксированным размером
	listScroll := container.NewScroll(cottageList)
	listScroll.Resize(fyne.NewSize(1000, 400))

	content := container.NewBorder(
		container.NewVBox(
			widget.NewCard("➕ Добавить новый домик", "", container.NewVBox(addForm, addBtn)),
			container.NewHBox(refreshBtn),
			widget.NewSeparator(),
		),
		nil, nil, nil,
		listScroll,
	)

	return container.NewTabItem("🏠 Домики", content)
}

// createTariffsTab создает вкладку тарифов с фиксированными размерами
func (a *StyledGuestApp) createTariffsTab() *container.TabItem {
	// Поиск
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("🔍 Поиск по названию...")

	// Список тарифов
	var filteredTariffs []models.Tariff
	var tariffList *widget.List

	// Статистика с фиксированным размером
	statsLabel := widget.NewLabel("")
	statsLabel.Resize(fyne.NewSize(200, 30))

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
			if searchEntry.Text != "" && !strings.Contains(strings.ToLower(tariff.Name), strings.ToLower(searchEntry.Text)) {
				continue
			}
			filteredTariffs = append(filteredTariffs, tariff)
		}

		if tariffList != nil {
			tariffList.Refresh()
		}

		// Обновляем статистику
		statsLabel.SetText(fmt.Sprintf("📊 Всего тарифов: %d", len(tariffs)))
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
			priceLabel.SetText(fmt.Sprintf("%.2f ₽/день", tariff.PricePerDay))

			editBtn.OnTapped = func() {
				a.showEditTariffDialogFixed(tariff)
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
	nameEntry.SetPlaceHolder("Название тарифа")

	priceEntry := widget.NewEntry()
	priceEntry.SetPlaceHolder("Цена за сутки")

	addForm := widget.NewForm(
		widget.NewFormItem("💰 Название", nameEntry),
		widget.NewFormItem("💵 Цена за сутки", priceEntry),
	)

	addBtn := widget.NewButtonWithIcon("Добавить тариф", theme.ContentAddIcon(), func() {
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
	addBtn.Importance = widget.HighImportance

	// Кнопка обновления с фиксированным размером
	refreshBtn := widget.NewButtonWithIcon("🔄 Обновить", theme.ViewRefreshIcon(), func() {
		updateTariffList()
		dialog.ShowInformation("Успешно", "Список обновлен", a.window)
	})
	refreshBtn.Resize(fyne.NewSize(150, 40))

	// Изначально загружаем тарифы
	updateTariffList()

	// Создаем контейнер со скроллом для списка с фиксированным размером
	listScroll := container.NewScroll(tariffList)
	listScroll.Resize(fyne.NewSize(1000, 400))

	content := container.NewBorder(
		container.NewVBox(
			searchEntry,
			widget.NewCard("➕ Добавить новый тариф", "", container.NewVBox(addForm, addBtn)),
			container.NewHBox(refreshBtn, statsLabel),
			widget.NewSeparator(),
		),
		nil, nil, nil,
		listScroll,
	)

	return container.NewTabItem("💰 Тарифы", content)
}

// updateStats обновляет статистику с ФИКСИРОВАННЫМИ размерами
func (a *StyledGuestApp) updateStats(label *widget.RichText) {
	// Получаем статистику за текущий месяц
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	bookings, err := a.bookingService.GetBookingsByDateRange(startOfMonth, endOfMonth)
	if err != nil {
		fyne.Do(func() {
			label.ParseMarkdown("❌ Ошибка загрузки статистики")
		})
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

	// Получаем статистику по домикам
	cottages, _ := a.cottageService.GetAllCottages()
	freeCottages := 0
	occupiedCottages := 0
	for _, c := range cottages {
		if c.Status == "free" {
			freeCottages++
		} else if c.Status == "occupied" {
			occupiedCottages++
		}
	}

	avgCheck := 0.0
	if activeBookings > 0 {
		avgCheck = totalRevenue / float64(activeBookings)
	}

	statsText := fmt.Sprintf(`## 📊 %s %d

**🏠 Домики:**
• Всего: %d
• 🟢 Свободно: %d  
• 🔴 Занято: %d

**📋 Бронирования:**
• Всего: %d
• Активных: %d

**💰 Доходы:**
• Общий: **%.0f ₽**
• Средний чек: **%.0f ₽**`,
		a.getMonthName(now.Month()), now.Year(),
		len(cottages), freeCottages, occupiedCottages,
		totalBookings, activeBookings,
		totalRevenue, avgCheck)

	fyne.Do(func() {
		label.ParseMarkdown(statsText)
		label.Resize(fyne.NewSize(260, 120)) // Updated width to match container
	})
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
	dialog.ShowInformation("🔧 В разработке", "Отчет о заполняемости будет реализован в следующей версии", a.window)
}

func (a *StyledGuestApp) generateRevenueReport() {
	dialog.ShowInformation("🔧 В разработке", "Отчет о доходах будет реализован в следующей версии", a.window)
}

func (a *StyledGuestApp) showUpcomingArrivals() {
	bookings, err := a.bookingService.GetUpcomingBookings()
	if err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	if len(bookings) == 0 {
		dialog.ShowInformation("📅 Информация", "Нет предстоящих заездов", a.window)
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

		// Создаем стилизованную карточку
		card := widget.NewCard(
			fmt.Sprintf("👤 %s • 🏠 %s", booking.GuestName, cottageName),
			fmt.Sprintf("📅 Заезд: %s", booking.CheckInDate.Format("02.01.2006")),
			container.NewVBox(
				widget.NewLabel(fmt.Sprintf("📞 %s", booking.Phone)),
				widget.NewLabel(fmt.Sprintf("💰 %.0f ₽", booking.TotalCost)),
			),
		)
		content.Add(card)
	}

	scroll := container.NewScroll(content)
	scroll.Resize(fyne.NewSize(500, 400))

	d := dialog.NewCustom("📅 Предстоящие заезды", "Закрыть", scroll, a.window)
	d.Resize(fyne.NewSize(600, 500))
	d.Show()
}

// showQuickBookingDialog показывает диалог быстрого бронирования (УНИКАЛЬНАЯ ВЕРСИЯ)
func (a *StyledGuestApp) showQuickBookingDialogFromSidePanel() {
	// Получаем список свободных домиков
	cottages, err := a.cottageService.GetFreeCottages()
	if err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	if len(cottages) == 0 {
		dialog.ShowInformation("ℹ️ Информация", "Нет свободных домиков", a.window)
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
		cottageOptions[i] = fmt.Sprintf("🏠 %s", c.Name)
	}
	cottageSelect := widget.NewSelect(cottageOptions, nil)

	tariffOptions := make([]string, len(tariffs))
	for i, t := range tariffs {
		tariffOptions[i] = fmt.Sprintf("💰 %s - %.0f ₽/сутки", t.Name, t.PricePerDay)
	}
	tariffSelect := widget.NewSelect(tariffOptions, nil)

	// Даты
	checkInDate := time.Now().Add(24 * time.Hour)
	checkOutDate := checkInDate.Add(24 * time.Hour)

	checkInPicker := ui.NewDatePickerButton("📅 Дата заезда", a.window, func(t time.Time) {
		checkInDate = t
	})
	checkInPicker.SetSelectedDate(checkInDate)

	checkOutPicker := ui.NewDatePickerButton("📅 Дата выезда", a.window, func(t time.Time) {
		checkOutDate = t
	})
	checkOutPicker.SetSelectedDate(checkOutDate)

	notesEntry := widget.NewMultiLineEntry()
	notesEntry.SetPlaceHolder("📝 Дополнительные примечания...")
	notesEntry.SetMinRowsVisible(3)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "👤 ФИО гостя", Widget: nameEntry},
			{Text: "📞 Телефон", Widget: phoneEntry},
			{Text: "📧 Email", Widget: emailEntry},
			{Text: "🏠 Домик", Widget: cottageSelect},
			{Text: "💰 Тариф", Widget: tariffSelect},
			{Text: "📅 Дата заезда", Widget: checkInPicker},
			{Text: "📅 Дата выезда", Widget: checkOutPicker},
			{Text: "📝 Примечания", Widget: notesEntry},
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
			dialog.ShowInformation("✅ Успешно", "Бронирование создано", a.window)
		},
		OnCancel: func() {},
	}

	d := dialog.NewCustom("✨ Новое бронирование", "Отмена", form, a.window)
	d.Resize(fyne.NewSize(500, 600))
	d.Show()
}

// showEditCottageDialogFixed показывает диалог редактирования домика (УНИКАЛЬНАЯ ВЕРСИЯ)
func (a *StyledGuestApp) showEditCottageDialogFixed(cottage models.Cottage) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(cottage.Name)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "🆔 ID", Widget: widget.NewLabel(fmt.Sprintf("%d", cottage.ID))},
			{Text: "🏠 Название", Widget: nameEntry},
			{Text: "📊 Статус", Widget: widget.NewLabel(cottage.Status)},
		},
		OnSubmit: func() {
			err := a.cottageService.UpdateCottageName(cottage.ID, nameEntry.Text)
			if err != nil {
				dialog.ShowError(err, a.window)
				return
			}
			dialog.ShowInformation("✅ Успешно", "Домик обновлен", a.window)
			a.calendarWidget.Update()
		},
	}

	d := dialog.NewCustom("✏️ Редактировать домик", "Отмена", form, a.window)
	d.Show()
}

// showEditTariffDialogFixed показывает диалог редактирования тарифа (УНИКАЛЬНАЯ ВЕРСИЯ)
func (a *StyledGuestApp) showEditTariffDialogFixed(tariff models.Tariff) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(tariff.Name)

	priceEntry := widget.NewEntry()
	priceEntry.SetText(fmt.Sprintf("%.2f", tariff.PricePerDay))

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "💰 Название", Widget: nameEntry},
			{Text: "💵 Цена за сутки", Widget: priceEntry},
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

			dialog.ShowInformation("✅ Успешно", "Тариф сохранен", a.window)
		},
	}

	d := dialog.NewCustom("✏️ Редактирование тарифа", "Отмена", form, a.window)
	d.Resize(fyne.NewSize(500, 400))
	d.Show()
}

// Вспомогательная функция для интерполяции цветов
func lerpColor(start, end color.Color, t float32) color.Color {
	r1, g1, b1, a1 := start.RGBA()
	r2, g2, b2, a2 := end.RGBA()

	return color.NRGBA{
		R: uint8(float32(r1>>8)*(1-t) + float32(r2>>8)*t),
		G: uint8(float32(g1>>8)*(1-t) + float32(g2>>8)*t),
		B: uint8(float32(b1>>8)*(1-t) + float32(b2>>8)*t),
		A: uint8(float32(a1>>8)*(1-t) + float32(a2>>8)*t),
	}
}

func (a *StyledGuestApp) showQuickBookingDialog() {
	// Получаем список свободных домиков
	cottages, err := a.cottageService.GetFreeCottages()
	if err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	if len(cottages) == 0 {
		dialog.ShowInformation("ℹ️ Информация", "Нет свободных домиков", a.window)
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
		cottageOptions[i] = fmt.Sprintf("🏠 %s", c.Name)
	}
	cottageSelect := widget.NewSelect(cottageOptions, nil)

	tariffOptions := make([]string, len(tariffs))
	for i, t := range tariffs {
		tariffOptions[i] = fmt.Sprintf("💰 %s - %.0f ₽/сутки", t.Name, t.PricePerDay)
	}
	tariffSelect := widget.NewSelect(tariffOptions, nil)

	// Даты
	checkInDate := time.Now().Add(24 * time.Hour)
	checkOutDate := checkInDate.Add(24 * time.Hour)

	checkInPicker := ui.NewDatePickerButton("📅 Дата заезда", a.window, func(t time.Time) {
		checkInDate = t
	})
	checkInPicker.SetSelectedDate(checkInDate)

	checkOutPicker := ui.NewDatePickerButton("📅 Дата выезда", a.window, func(t time.Time) {
		checkOutDate = t
	})
	checkOutPicker.SetSelectedDate(checkOutDate)

	notesEntry := widget.NewMultiLineEntry()
	notesEntry.SetPlaceHolder("📝 Дополнительные примечания...")
	notesEntry.SetMinRowsVisible(3)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "👤 ФИО гостя", Widget: nameEntry},
			{Text: "📞 Телефон", Widget: phoneEntry},
			{Text: "📧 Email", Widget: emailEntry},
			{Text: "🏠 Домик", Widget: cottageSelect},
			{Text: "💰 Тариф", Widget: tariffSelect},
			{Text: "📅 Дата заезда", Widget: checkInPicker},
			{Text: "📅 Дата выезда", Widget: checkOutPicker},
			{Text: "📝 Примечания", Widget: notesEntry},
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
			dialog.ShowInformation("✅ Успешно", "Бронирование создано", a.window)
		},
		OnCancel: func() {},
	}

	d := dialog.NewCustom("✨ Новое бронирование", "Отмена", form, a.window)
	d.Resize(fyne.NewSize(500, 600))
	d.Show()
}

// showEditCottageDialog - основная версия диалога редактирования домика
func (a *StyledGuestApp) showEditCottageDialog(cottage models.Cottage) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(cottage.Name)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "🆔 ID", Widget: widget.NewLabel(fmt.Sprintf("%d", cottage.ID))},
			{Text: "🏠 Название", Widget: nameEntry},
			{Text: "📊 Статус", Widget: widget.NewLabel(cottage.Status)},
		},
		OnSubmit: func() {
			err := a.cottageService.UpdateCottageName(cottage.ID, nameEntry.Text)
			if err != nil {
				dialog.ShowError(err, a.window)
				return
			}
			dialog.ShowInformation("✅ Успешно", "Домик обновлен", a.window)
			a.calendarWidget.Update()
		},
	}

	d := dialog.NewCustom("✏️ Редактировать домик", "Отмена", form, a.window)
	d.Show()
}

// showEditTariffDialog - основная версия диалога редактирования тарифа
func (a *StyledGuestApp) showEditTariffDialog(tariff models.Tariff) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(tariff.Name)

	priceEntry := widget.NewEntry()
	priceEntry.SetText(fmt.Sprintf("%.2f", tariff.PricePerDay))

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "💰 Название", Widget: nameEntry},
			{Text: "💵 Цена за сутки", Widget: priceEntry},
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

			dialog.ShowInformation("✅ Успешно", "Тариф сохранен", a.window)
		},
	}

	d := dialog.NewCustom("✏️ Редактирование тарифа", "Отмена", form, a.window)
	d.Resize(fyne.NewSize(500, 400))
	d.Show()
}
