package ui

import (
	"fmt"
	"guest-management/internal/models"
	"guest-management/internal/service"
	"image/color"
	"log"
	"os/exec"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// GuestApp — основное приложение для учёта гостей
// содержит сервисы для гостей и домиков
type GuestApp struct {
	app            fyne.App
	window         fyne.Window
	guestService   *service.GuestService
	cottageService *service.CottageService
}

// NewGuestApp создаёт новое приложение
func NewGuestApp(guestService *service.GuestService, cottageService *service.CottageService) *GuestApp {
	a := app.New()
	w := a.NewWindow("Учет гостей - База отдыха")
	w.Resize(fyne.NewSize(800, 600))

	return &GuestApp{app: a, window: w, guestService: guestService, cottageService: cottageService}
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
	)
	a.window.SetContent(tabs)
}

// createGuestTab создаёт вкладку для регистрации гостей
func (a *GuestApp) createGuestTab() *container.TabItem {
	// Поля формы
	nameEntry := widget.NewEntry()
	emailEntry := widget.NewEntry()
	phoneEntry := widget.NewEntry()
	cottageSelect := widget.NewSelect([]string{}, nil)
	documentPath := widget.NewLabel("")

	// Функция обновления списка свободных домиков
	updateCottages := func() {
		cottages, err := a.cottageService.GetFreeCottages()
		if err != nil {
			log.Println("Error fetching cottages:", err)
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
	updateCottages() // Первоначальная загрузка

	// Кнопка выбора документа
	fileButton := widget.NewButton("Выбрать документ", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, a.window)
				return
			}
			if reader == nil {
				return // Пользователь отменил выбор
			}
			documentPath.SetText(reader.URI().Path())
		}, a.window)
	})

	// Форма регистрации
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ФИО", Widget: nameEntry},
			{Text: "Email", Widget: emailEntry},
			{Text: "Телефон", Widget: phoneEntry},
			{Text: "Домик", Widget: cottageSelect},
			{Text: "Документ", Widget: container.NewHBox(
				fileButton,
				documentPath,
			)},
		},
		OnSubmit: func() {
			// Валидация полей
			if nameEntry.Text == "" ||
				emailEntry.Text == "" ||
				phoneEntry.Text == "" ||
				cottageSelect.Selected == "" {
				dialog.ShowError(fmt.Errorf("Все поля обязательны для заполнения"), a.window)
				return
			}

			// Поиск ID выбранного домика
			cottages, err := a.cottageService.GetFreeCottages()
			if err != nil {
				dialog.ShowError(fmt.Errorf("Ошибка получения данных домиков"), a.window)
				return
			}

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
				dialog.ShowError(fmt.Errorf("Выбранный домик больше не доступен"), a.window)
				updateCottages()
				return
			}

			// Создание объекта гостя
			guest := models.Guest{
				FullName:         nameEntry.Text,
				Email:            emailEntry.Text,
				Phone:            phoneEntry.Text,
				DocumentScanPath: documentPath.Text,
			}

			// Регистрация гостя
			if err := a.guestService.RegisterGuest(guest, cottageID); err != nil {
				dialog.ShowError(fmt.Errorf("Ошибка регистрации: %v", err), a.window)
				return
			}

			// Очистка формы
			nameEntry.SetText("")
			emailEntry.SetText("")
			phoneEntry.SetText("")
			cottageSelect.ClearSelected()
			documentPath.SetText("")

			// Обновление списка домиков
			updateCottages()

			dialog.ShowInformation("Успешно!", "Гость зарегистрирован", a.window)
		},
	}

	return container.NewTabItem("Регистрация гостей",
		container.NewVScroll(
			container.NewPadded(form),
		),
	)
}

// createCottagesTab создаёт вкладку для управления домиками
func (a *GuestApp) createCottagesTab() *container.TabItem {
	content := container.NewGridWithColumns(5)

	updateContent := func() {
		content.RemoveAll()
		cottages, err := a.cottageService.GetAllCottages()
		if err != nil {
			log.Println("Error fetching cottages:", err)
			return
		}

		for _, c := range cottages {
			btn := NewCottageButton(c, func() {
				if c.Status == "occupied" {
					a.showCottageDetails(c)
				}
			})
			content.Add(btn)
		}
		content.Refresh()
	}
	updateContent()

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

// showCottageDetails отображает информацию о госте и кнопку "Выселить"
func (a *GuestApp) showCottageDetails(c models.Cottage) {
	guest, err := a.guestService.GetGuestByCottageID(c.ID)
	if err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	detailsWindow := a.app.NewWindow(fmt.Sprintf("Домик: %s", c.Name))
	content := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("ФИО: %s", guest.FullName)),
		widget.NewLabel(fmt.Sprintf("Email: %s", guest.Email)),
		widget.NewLabel(fmt.Sprintf("Телефон: %s", guest.Phone)),
		widget.NewLabel(fmt.Sprintf("Документ: %s", guest.DocumentScanPath)),
		widget.NewButton("Выселить", func() {
			if err := a.guestService.CheckOutGuest(c.ID); err != nil {
				dialog.ShowError(err, detailsWindow)
				return
			}
			detailsWindow.Close()
			a.createUI() // Обновляем интерфейс
		}),
		// Добавляем новую кнопку
		widget.NewButton("Открыть папку с документами", func() {
			if guest.DocumentScanPath == "" {
				dialog.ShowInformation("Информация", "Документы отсутствуют", detailsWindow)
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
				dialog.ShowError(fmt.Errorf("Неподдерживаемая ОС"), detailsWindow)
			}
		}),
	)

	detailsWindow.SetContent(content)
	detailsWindow.Show()
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
// реализует fyne.Tappable и имеет свой рендерер
// -----------------------------------------------
// структура кнопки
// -----------------------------------------------
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
	return &cottageButtonRenderer{box: box, objects: []fyne.CanvasObject{box}}
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
// -----------------------------------------------
type cottageButtonRenderer struct {
	box     *fyne.Container
	objects []fyne.CanvasObject
}

func (r *cottageButtonRenderer) Layout(size fyne.Size) {
	r.box.Resize(size)
}
func (r *cottageButtonRenderer) MinSize() fyne.Size {
	return r.box.MinSize()
}
func (r *cottageButtonRenderer) Refresh() {
	canvas.Refresh(r.box)
}
func (r *cottageButtonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}
func (r *cottageButtonRenderer) Destroy() {}
