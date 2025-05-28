package app

import (
	"fmt"
	"log"
	"strconv"
	"strings"
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

	// Функция обновления списка домиков
	var updateCottageList func()

	// Фильтр по статусу
	statusFilter := widget.NewSelect([]string{"Все", "Свободные", "Занятые"}, func(s string) {
		updateCottageList()
	})
	statusFilter.PlaceHolder = "Все"

	// Форма поиска
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Поиск по названию...")
	searchEntry.OnChanged = func() {
		updateCottageList()
	}

	// Контейнер для списка домиков
	cottagesContainer := container.NewVBox()
	cottagesContainer.Objects = make([]fyne.CanvasObject, 0)

	// Функция обновления списка домиков
	updateCottageList = func() {
		cottages, err := a.cottageService.GetAllCottages()
		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}

		filteredCottages := make([]models.Cottage, 0)
		for _, cottage := range cottages {
			if statusFilter.Selected == "Свободные" && cottage.Status != "free" {
				continue
			}
			if statusFilter.Selected == "Занятые" && cottage.Status != "occupied" {
				continue
			}
			if searchEntry.Text != "" && !strings.Contains(strings.ToLower(cottage.Name), strings.ToLower(searchEntry.Text)) {
				continue
			}
			filteredCottages = append(filteredCottages, cottage)
		}

		// Очищаем контейнер
		cottagesContainer.Objects = cottagesContainer.Objects[:0]

		// Добавляем домики
		for _, cottage := range filteredCottages {
			// Контейнер для каждого домика
			cottageContainer := container.NewVBox(
				container.NewCard(
					cottage.Name,
					fmt.Sprintf("Статус: %s", cottage.Status),
					container.NewVBox(
						widget.NewButton("Редактировать", func() {
							a.showEditCottageDialog(cottage, updateCottageList)
						}),
						widget.NewButton("Забронировать", func() {
							bookingDialog := dialog.NewCustomConfirm("Бронирование", "Забронировать", "Отмена", func(b bool) {
								if b {
									// TODO: Implement booking functionality
								}
							}, a.window)
							bookingDialog.SetOnClosed(func() {
								bookingDialog = nil
							})
							bookingDialog.Show()
						}),
					),
				),
			)

			// Кнопка удаления
			deleteBtn := widget.NewButton("Удалить", func() {
				if cottage.Status == "occupied" {
					dialog.ShowError(
						fmt.Errorf("невозможно удалить занятый домик"),
						a.window,
					)
					return
				}

				dialog.ShowConfirm("Подтверждение",
					fmt.Sprintf("Вы уверены, что хотите удалить домик '%s' (ID: %d)?", cottage.Name, cottage.ID),
					func(ok bool) {
						if ok {
							err := a.cottageService.DeleteCottage(cottage.ID)
							if err != nil {
								dialog.ShowError(err, a.window)
								return
							}
							for i, c := range a.cottages {
								if c.ID == cottage.ID {
									a.cottages = append(a.cottages[:i], a.cottages[i+1:]...)
									break
								}
							}
							updateCottageList()
							dialog.ShowInformation("Успешно", fmt.Sprintf("Домик '%s' удален", cottage.Name), a.window)
						}
					},
					a.window,
				)
			})
			cottageContainer.Add(deleteBtn)

			// Booking button with proper state handling
			bookingBtn := widget.NewButton("Заселить", func() {
				bookingDialog := dialog.NewCustomConfirm("Заселение", "Заселить", "Отмена",
					func(b bool) {
						if b {
							// TODO: Implement booking functionality
						}
					},
					a.window)
				bookingDialog.SetOnClosed(func() {
					bookingDialog = nil
				})
				bookingDialog.Show()
				updateCottageList()
			})

			cottageContainer.Add(deleteBtn)
			cottageContainer.Add(bookingBtn)

    // Кнопка добавления
    addBtn := widget.NewButton("Добавить", func() {
        // Создаем диалог добавления
        dialog.ShowCustom("Добавление домика", "Добавить", container.NewVBox(
            nameEntry,
            widget.NewButton("Добавить", func() {
                if nameEntry.Text == "" {
                    return
                }

                c, err := a.cottageService.CreateCottage(nameEntry.Text)
                if err != nil {
                    log.Printf("Ошибка при создании домика: %v", err)
                    return
                }

                a.cottages = append(a.cottages, *c)
                updateCottageList()
                dialog.ShowInformation("Успешно", "Список обновлен", a.window)
            }),
        ))
    })

    // Кнопка обновления
    refreshBtn := widget.NewButton("Обновить", func() {
        updateCottageList()
    })

    // Кнопка бронирования
    bookingBtn := widget.NewButton("Заселить", func() {
        bookingDialog := dialog.NewCustomConfirm("Заселение", "Заселить", "Отмена",
            func(b bool) {
                if b {
                    // TODO: Implement booking functionality
                }
            },
            a.window)
        bookingDialog.SetOnClosed(func() {
            bookingDialog = nil
        })
        bookingDialog.Show()
        updateCottageList()
    })

    // Добавляем кнопки в контейнер
    btnContainer := container.NewHBox(addBtn, bookingBtn, refreshBtn)
		}
		statsLabel.SetText(fmt.Sprintf("Всего домиков: %d | Свободно: %d | Занято: %d",
			total, free, occupied))
	}
	updateStats()

	// Создаем контейнер со скроллом для списка
	listScroll := container.NewScroll(cottagesContainer)
	listScroll.SetMinSize(fyne.NewSize(600, 400))

	content := container.NewBorder(
		container.NewVBox(
			container.NewHBox(refreshBtn, statsLabel),
			widget.NewSeparator(),
		),
		nil, nil, nil,
		listScroll,
	)

	return container.NewTabItem("Домики", content)
}

// showEditCottageDialog показывает диалог редактирования домика
func (a *GuestApp) showEditCottageDialog(cottage models.Cottage, onUpdate func()) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(cottage.Name)

	form := &widget.Form{
			Items: []*widget.FormItem{
				{Text: "ID домика", Widget: widget.NewLabel(fmt.Sprintf("%d", cottage.ID))},
				{Text: "Название", Widget: nameEntry},
				{Text: "Статус", Widget: widget.NewLabel(func() string {
					if cottage.Status == "occupied" {
						return "Занят"
					}
					return "Свободен"
				}())},
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
			err := a.cottageService.UpdateCottageName(cottage.ID, nameEntry.Text)
			if err != nil {
				dialog.ShowError(err, a.window)
				return
			}

			onUpdate()

			// Обновляем календарь после изменения названия
			if a.calendarWidget != nil {
				a.calendarWidget.Update()
			}

			dialog.ShowInformation("Успешно", "Название домика обновлено", a.window)
		},
	}

	dialog.ShowCustom("Редактировать домик", "Отмена", form, a.window)
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
				dialog.ShowConfirm("Удаление тарифа", "Вы уверены, что хотите удалить этот тариф?", func(ok bool) {
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
		}
	)

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

	form := &widget.Form{
	Items: []*widget.FormItem{
		{Text: "Название тарифа", Widget: nameEntry},
		{Text: "Цена за день (руб.)", Widget: priceEntry},
	},
	OnSubmit: func() {
		if nameEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("название не может быть пустым"), a.window)
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

		onUpdate()
		dialog.ShowInformation("Успешно", "Тариф обновлен", a.window)
	},

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

// ...
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
