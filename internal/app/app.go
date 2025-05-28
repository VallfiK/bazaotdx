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
	cottages             []models.Cottage
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
	// Создаем виджеты календаря и списка бронирований
	calendarWidget := widget.NewList(
		func() int {
			// TODO: Implement calendar widget
			return 0
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Loading...")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			// TODO: Implement calendar item rendering
		},
	)

	bookingListWidget := a.createBookingList()

	// Устанавливаем взаимные обновления
	bookingListWidget.OnSelected = func(id widget.ListItemID) {
		// TODO: Implement refresh logic
	}

	// Создаем вкладки (убираем ненужные)
	tabs := container.NewAppTabs(
		container.NewTabItem("Календарь", container.NewHSplit(calendarWidget, bookingListWidget)),
		a.createTariffsTab(),
		a.createReportsTab(),
	)

	// Устанавливаем основное содержимое окна
	a.window.SetContent(tabs)

	// Устанавливаем иконки для вкладок
	tabs.Items[0].Icon = theme.CalendarIcon()
	tabs.Items[1].Icon = theme.ListIcon()
	tabs.Items[2].Icon = theme.SettingsIcon()
	tabs.Items[3].Icon = theme.DocumentIcon()
}

// createBookingList создает виджет списка бронирований
func (a *GuestApp) createBookingList() *widget.List {
	// Создаем список бронирований
	list := widget.NewList(
		func() int {
			bookings, err := a.bookingService.GetBookingsByStatus(models.BookingStatusBooked, time.Now())
			if err != nil {
				log.Printf("Error loading bookings: %v", err)
				return 0
			}
			return len(bookings)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Loading...")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			bookings, err := a.bookingService.GetBookingsByStatus(models.BookingStatusBooked, time.Now())
			if err != nil {
				log.Printf("Error loading bookings: %v", err)
				return
			}
			if id < 0 || id >= widget.ListItemID(len(bookings)) {
				return
			}
			booking := bookings[id]
			label := item.(*widget.Label)
			label.SetText(fmt.Sprintf("%s - %s", booking.GuestName, booking.Status))
		},
	)

	// Добавляем кнопки управления
	list.OnSelected = func(id widget.ListItemID) {
		bookings, err := a.bookingService.GetBookingsByStatus(models.BookingStatusBooked, time.Now())
		if err != nil {
			log.Printf("Error loading bookings: %v", err)
			return
		}
		if id < 0 || id >= widget.ListItemID(len(bookings)) {
			return
		}
		booking := bookings[id]
		a.showBookingActions(booking)
	}

	return list
}

func (a *GuestApp) showBookingActions(booking models.Booking) {
	// Находим название домика
	cottageName := ""
	for _, c := range a.cottages {
		if c.ID == booking.CottageID {
			cottageName = c.Name
			break
		}
	}

	// Создаем контент для диалога
	content := container.NewVBox(
		widget.NewCard("Информация о брони", "",
			container.NewVBox(
				widget.NewLabel(fmt.Sprintf("Гость: %s", booking.GuestName)),
				widget.NewLabel(fmt.Sprintf("Телефон: %s", booking.Phone)),
				widget.NewLabel(fmt.Sprintf("Email: %s", booking.Email)),
				widget.NewLabel(fmt.Sprintf("Домик: %s", cottageName)),
				widget.NewLabel(fmt.Sprintf("Заезд: %s", booking.CheckInDate.Format("02.01.2006 15:04"))),
				widget.NewLabel(fmt.Sprintf("Выезд: %s", booking.CheckOutDate.Format("02.01.2006 15:04"))),
				widget.NewLabel(fmt.Sprintf("Статус: %s", booking.Status)),
				widget.NewLabel(fmt.Sprintf("Стоимость: %.2f руб.", booking.TotalCost)),
				widget.NewLabel(fmt.Sprintf("Создано: %s", booking.CreatedAt.Format("02.01.2006 15:04"))),
			),
		),
	)

	// Кнопки действий
	actions := container.NewHBox()

	switch booking.Status {
	case models.BookingStatusBooked:
		actions.Add(widget.NewButton("Заселить", func() {
			dialog.ShowConfirm("Подтверждение", "Заселить гостя?", func(ok bool) {
				if ok {
					err := a.bookingService.CheckInBooking(booking.ID)
					if err != nil {
						dialog.ShowError(err, a.window)
						return
					}
					a.createUI()
					dialog.ShowInformation("Успешно", "Гость заселен", a.window)
				}
			}, a.window)
		}))

		actions.Add(widget.NewButton("Отменить", func() {
			dialog.ShowConfirm("Подтверждение", "Отменить бронирование?", func(ok bool) {
				if ok {
					err := a.bookingService.CancelBooking(booking.ID)
					if err != nil {
						dialog.ShowError(err, a.window)
						return
					}
					a.createUI()
					dialog.ShowInformation("Успешно", "Бронирование отменено", a.window)
				}
			}, a.window)
		}))

	case models.BookingStatusCheckedIn:
		// Добавьте действия для заселенных гостей
	}

	content.Add(actions)

	d := dialog.NewCustom("Действия с бронированием", "Закрыть", content, a.window)
	d.Resize(fyne.NewSize(500, 500))
	d.Show()
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

	return container.NewTabItem("Тарифы", container.NewVBox(form, tariffList))
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
	
	// Создаем контейнер с кнопками
	content := container.NewVBox(
		widget.NewLabel("Выберите тип отчета:"),
		container.NewHBox(occupancyBtn, revenueBtn),
		upcomingBtn,
	)
	
	return container.NewTabItem("Отчеты", content)
}

func (a *GuestApp) generateOccupancyReport() {
	// TODO: Реализовать генерацию отчета по заполняемости
}

func (a *GuestApp) generateRevenueReport() {
	// TODO: Реализовать генерацию отчета по доходам
}

func (a *GuestApp) showUpcomingArrivals() {
	// TODO: Реализовать просмотр предстоящих заездов
}
