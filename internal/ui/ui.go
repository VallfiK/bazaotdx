package ui

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/VallfIK/bazaotdx/internal/models"
	"github.com/VallfIK/bazaotdx/internal/service"
)

// GuestApp — основное приложение для учёта гостей
// содержит сервисы для гостей и домиков
type GuestApp struct {
	app                   fyne.App
	window                fyne.Window
	guestService          *service.GuestService
	cottageService        *service.CottageService
	tariffService         *service.TariffService // Добавлено
	updateCottagesContent func()
}

// NewGuestApp создаёт новое приложение
func NewGuestApp(
	guestService *service.GuestService,
	cottageService *service.CottageService,
	tariffService *service.TariffService, // Добавлено
) *GuestApp {
	a := app.New()
	w := a.NewWindow("Учет гостей - База отдыха")
	w.Resize(fyne.NewSize(800, 600))

	return &GuestApp{
		app:            a,
		window:         w,
		guestService:   guestService,
		cottageService: cottageService,
		tariffService:  tariffService, // Добавлено
	}
}

// Run запускает приложение
func (a *GuestApp) Run() {
	a.createUI()
	a.window.ShowAndRun()
}

// createUI формирует основное содержание окна с вкладками
func (a *GuestApp) createUI() {
	tabs := container.NewAppTabs(
		a.createGuestTab(),
		a.createCottagesTab(),
		a.createTariffsTab(),
	)
	a.window.SetContent(tabs)
}

func (a *GuestApp) createTariffsTab() *container.TabItem {
	nameEntry := widget.NewEntry()
	priceEntry := widget.NewEntry()
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Название тарифа", Widget: nameEntry},
			{Text: "Цена за день", Widget: priceEntry},
		},
		OnSubmit: func() {
			price, _ := strconv.ParseFloat(priceEntry.Text, 64)
			err := a.tariffService.CreateTariff(nameEntry.Text, price)
			if err != nil {
				dialog.ShowError(err, a.window)
				return
			}
			dialog.ShowInformation("Успех", "Тариф создан", a.window)
		},
	}
	return container.NewTabItem("Тарифы", form)
}

func calculateDays(checkIn, checkOut time.Time) float64 {
	checkIn = checkIn.Add(-14 * time.Hour) // Нормализация до 00:00
	checkOut = checkOut.Add(-12 * time.Hour)

	days := checkOut.Sub(checkIn).Hours() / 24
	return math.Ceil(days*100) / 100 // Округление до 2 знаков
}

// createGuestTab создаёт вкладку для регистрации гостей
func (a *GuestApp) createGuestTab() *container.TabItem {
	// Поля формы
	nameEntry := widget.NewEntry()
	emailEntry := widget.NewEntry()
	phoneEntry := widget.NewEntry()
	cottageSelect := widget.NewSelect([]string{}, nil)
	documentPath := widget.NewLabel("")
	tariffSelect := widget.NewSelect([]string{}, nil)
	costLabel := widget.NewLabel("Стоимость: 0.00 руб.")

	// Переменные состояния
	var (
		checkIn  time.Time
		checkOut time.Time
		cottages []models.Cottage
		tariffs  []models.Tariff
	)

	// Текстовые поля для ввода дат
	checkInEntry := widget.NewEntry()
	checkInEntry.SetPlaceHolder("дд.мм.гггг 14:00")
	checkOutEntry := widget.NewEntry()
	checkOutEntry.SetPlaceHolder("дд.мм.гггг 12:00")

	// Функция обработки ввода дат
	setupDateEntry := func(entry *widget.Entry, setTime *time.Time, hour int) {
		entry.OnChanged = func(s string) {
			t, err := time.Parse("02.01.2006 15:04", s)
			if err == nil {
				*setTime = time.Date(
					t.Year(), t.Month(), t.Day(),
					hour, 0, 0, 0, time.Local,
				)
				updateCost()
			}
		}
	}

	// Инициализация обработчиков дат
	setupDateEntry(checkInEntry, &checkIn, 14)
	setupDateEntry(checkOutEntry, &checkOut, 12)

	// Инициализация тарифов
	initTariffs := func() {
		var err error
		tariffs, err = a.tariffService.GetTariffs()
		if err != nil {
			dialog.ShowError(fmt.Errorf("Ошибка загрузки тарифов"), a.window)
			return
		}

		tariffSelect.Options = make([]string, len(tariffs))
		for i, t := range tariffs {
			tariffSelect.Options[i] = fmt.Sprintf("%s (%.2f руб./день)", t.Name, t.PricePerDay)
		}
		tariffSelect.Refresh()
	}
	initTariffs()

	// Обновление списка домиков
	updateCottages := func() {
		var err error
		cottages, err = a.cottageService.GetFreeCottages()
		if err != nil {
			dialog.ShowError(fmt.Errorf("Ошибка загрузки домиков"), a.window)
			return
		}

		options := make([]string, len(cottages))
		for i, c := range cottages {
			options[i] = c.Name
		}
		cottageSelect.Options = options
		cottageSelect.Refresh()
	}
	updateCottages()

	// Расчет стоимости (теперь объявлена до использования)
	updateCost := func() {
		if checkIn.IsZero() || checkOut.IsZero() || tariffSelect.SelectedIndex() < 0 {
			costLabel.SetText("Стоимость: -")
			return
		}

		if checkOut.Before(checkIn) {
			costLabel.SetText("Дата выезда раньше заезда!")
			return
		}

		days := checkOut.Sub(checkIn).Hours() / 24
		price := tariffs[tariffSelect.SelectedIndex()].PricePerDay
		cost := days * price
		costLabel.SetText(fmt.Sprintf("Стоимость: %.2f руб.", cost))
	}

	// Добавляем обработчик изменения тарифа
	tariffSelect.OnChanged = func(_ string) {
		updateCost()
	}

	// Выбор документа
	fileButton := widget.NewButton("Выбрать документ", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				documentPath.SetText(reader.URI().Path())
			}
		}, a.window)
	})

	// Форма регистрации
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ФИО", Widget: nameEntry},
			{Text: "Email", Widget: emailEntry},
			{Text: "Телефон", Widget: phoneEntry},
			{Text: "Домик", Widget: cottageSelect},
			{Text: "Дата заезда", Widget: checkInEntry},
			{Text: "Дата выезда", Widget: checkOutEntry},
			{Text: "Тариф", Widget: tariffSelect},
			{Text: "Документ", Widget: container.NewHBox(fileButton, documentPath)},
			{Text: "Стоимость", Widget: costLabel},
		},
		OnSubmit: func() {
			// Валидация
			if nameEntry.Text == "" || emailEntry.Text == "" || phoneEntry.Text == "" ||
				cottageSelect.Selected == "" || checkIn.IsZero() || checkOut.IsZero() ||
				tariffSelect.SelectedIndex() < 0 {
				dialog.ShowError(fmt.Errorf("Заполните все обязательные поля"), a.window)
				return
			}

			if checkOut.Before(checkIn) {
				dialog.ShowError(fmt.Errorf("Дата выезда раньше заезда"), a.window)
				return
			}

			// Поиск ID домика
			var cottageID int
			found := false
			for _, c := range cottages {
				if c.Name == cottageSelect.Selected {
					cottageID = c.ID
					found = true
					break
				}
			}

			if !found {
				dialog.ShowError(fmt.Errorf("Домик недоступен"), a.window)
				updateCottages()
				return
			}

			// Регистрация гостя
			guest := models.Guest{
				FullName:         nameEntry.Text,
				Email:            emailEntry.Text,
				Phone:            phoneEntry.Text,
				DocumentScanPath: documentPath.Text,
				CheckInDate:      checkIn,
				CheckOutDate:     checkOut,
				TariffID:         tariffs[tariffSelect.SelectedIndex()].ID,
			}

			if err := a.guestService.RegisterGuest(guest, cottageID); err != nil {
				dialog.ShowError(fmt.Errorf("Ошибка регистрации: %v", err), a.window)
				return
			}

			// Сброс формы
			nameEntry.SetText("")
			emailEntry.SetText("")
			phoneEntry.SetText("")
			cottageSelect.ClearSelected()
			documentPath.SetText("")
			checkInEntry.SetText("")
			checkOutEntry.SetText("")
			tariffSelect.ClearSelected()
			costLabel.SetText("Стоимость: 0.00 руб.")
			updateCottages()

			dialog.ShowInformation("Успех", "Регистрация завершена", a.window)
		},
	}

	return container.NewTabItem("Регистрация гостей",
		container.NewVScroll(
			container.NewPadded(form),
		))
}

func updateCost() {
	panic("unimplemented")
}

// createCottagesTab создаёт вкладку для управления домиками
func (a *GuestApp) createCottagesTab() *container.TabItem {
	content := container.NewGridWithColumns(5)

	// ИСПРАВЛЕНИЕ: Объявляем переменную для функции обновления
	var updateContent func()

	// Определяем функцию обновления
	updateContent = func() {
		content.RemoveAll()
		cottages, err := a.cottageService.GetAllCottages()
		if err != nil {
			log.Println("Error fetching cottages:", err)
			return
		}

		for _, c := range cottages {
			btn := NewCottageButton(c, func(cottage models.Cottage) func() {
				return func() {
					if cottage.Status == "occupied" {
						a.showCottageDetails(cottage, updateContent)
					}
				}
			}(c)) // Замыкание для передачи правильного cottage
			content.Add(btn)
		}
		content.Refresh()
	}

	// Сохраняем функцию обновления в приложении
	a.updateCottagesContent = updateContent

	updateContent() // Первоначальная загрузка

	addButton := widget.NewButton("Добавить домик", func() {
		dialog.ShowEntryDialog("Новый домик", "Введите название домика:", func(name string) {
			if name != "" {
				if err := a.cottageService.AddCottage(name); err != nil {
					dialog.ShowError(err, a.window)
					return
				}
				updateContent()
			}
		}, a.window)
	})

	return container.NewTabItem("Домики",
		container.NewBorder(addButton, nil, nil, nil, content))
}

// ИСПРАВЛЕНИЕ: Отображаем информацию о домике в модальном окне по центру
func (a *GuestApp) showCottageDetails(c models.Cottage, refreshCallback func()) {
	guest, err := a.guestService.GetGuestByCottageID(c.ID)
	if err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	// Создаем контент для модального окна
	content := container.NewVBox(
		widget.NewCard("", fmt.Sprintf("Информация о домике: %s", c.Name),
			container.NewVBox(
				widget.NewLabel(fmt.Sprintf("ФИО: %s", guest.FullName)),
				widget.NewLabel(fmt.Sprintf("Email: %s", guest.Email)),
				widget.NewLabel(fmt.Sprintf("Телефон: %s", guest.Phone)),
				widget.NewLabel(fmt.Sprintf("Документ: %s", guest.DocumentScanPath)),
				container.NewHBox(
					widget.NewButton("Выселить", func() {
						// Подтверждение выселения
						dialog.ShowConfirm("Подтверждение",
							fmt.Sprintf("Вы уверены, что хотите выселить %s из домика %s?", guest.FullName, c.Name),
							func(confirm bool) {
								if confirm {
									if err := a.guestService.CheckOutGuest(c.ID); err != nil {
										dialog.ShowError(err, a.window)
										return
									}
									// Вызываем колбэк для обновления
									if refreshCallback != nil {
										refreshCallback()
									}
									dialog.ShowInformation("Успешно!", "Гость выселен", a.window)
								}
							}, a.window)
					}),
					widget.NewButton("Открыть папку", func() {
						if guest.DocumentScanPath == "" {
							dialog.ShowInformation("Информация", "Документы отсутствуют", a.window)
							return
						}

						// Получаем путь к папке
						dirPath := filepath.Dir(guest.DocumentScanPath)

						// Открываем папку в файловом менеджере
						switch runtime.GOOS {
						case "windows":
							exec.Command("explorer.exe", dirPath).Start()
						case "darwin":
							exec.Command("open", dirPath).Start()
						case "linux":
							exec.Command("xdg-open", dirPath).Start()
						default:
							dialog.ShowError(fmt.Errorf("Неподдерживаемая ОС"), a.window)
						}
					}),
				),
			),
		),
	)

	// Создаем и показываем модальное окно
	modal := dialog.NewCustom(
		"", // Без заголовка, так как у нас есть Card с заголовком
		"Закрыть",
		content,
		a.window,
	)

	// Устанавливаем фиксированный размер модального окна
	modal.Resize(fyne.NewSize(400, 300))
	modal.Show()
}

// colorForStatus возвращает цвет для фона в зависимости от статуса
func colorForStatus(status string) color.Color {
	if status == "free" {
		return color.NRGBA{R: 0, G: 200, B: 0, A: 255}
	}
	return color.NRGBA{R: 200, G: 0, B: 0, A: 255}
}

// textColorForStatus возвращает цвет текста (белый)
func textColorForStatus(status string) color.Color {
	return color.White
}

// CottageButton — кастомная кнопка для отображения статуса домика
type CottageButton struct {
	widget.BaseWidget
	Cottage  models.Cottage
	OnTapped func()
}

// NewCottageButton создаёт новую кнопку
func NewCottageButton(c models.Cottage, tapped func()) *CottageButton {
	btn := &CottageButton{Cottage: c, OnTapped: tapped}
	btn.ExtendBaseWidget(btn)
	return btn
}

// CreateRenderer возвращает рендерер для кнопки
func (b *CottageButton) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(colorForStatus(b.Cottage.Status))
	txt := canvas.NewText(b.Cottage.Name, textColorForStatus(b.Cottage.Status))
	txt.Alignment = fyne.TextAlignCenter
	box := container.NewMax(bg, container.NewCenter(txt))
	return &cottageButtonRenderer{
		bg:      bg,
		txt:     txt,
		box:     box,
		objects: []fyne.CanvasObject{box},
		cottage: &b.Cottage, // Сохраняем ссылку на cottage для обновления
	}
}

// Tapped обрабатывает нажатие
func (b *CottageButton) Tapped(_ *fyne.PointEvent) {
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

// TappedSecondary для правого клика (не используется)
func (b *CottageButton) TappedSecondary(_ *fyne.PointEvent) {}

// cottageButtonRenderer рендерер для CottageButton
type cottageButtonRenderer struct {
	bg      *canvas.Rectangle
	txt     *canvas.Text
	box     *fyne.Container
	objects []fyne.CanvasObject
	cottage *models.Cottage
}

func (r *cottageButtonRenderer) Layout(size fyne.Size) {
	r.box.Resize(size)
}

func (r *cottageButtonRenderer) MinSize() fyne.Size {
	return r.box.MinSize()
}

// ИСПРАВЛЕНИЕ: Обновляем цвета при рефреше
func (r *cottageButtonRenderer) Refresh() {
	if r.cottage != nil {
		r.bg.FillColor = colorForStatus(r.cottage.Status)
		r.txt.Color = textColorForStatus(r.cottage.Status)
	}
	canvas.Refresh(r.box)
}

func (r *cottageButtonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *cottageButtonRenderer) Destroy() {}
