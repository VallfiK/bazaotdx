// internal/ui/styled_widgets.go
package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// StyledCard создает стилизованную карточку с тенью
func StyledCard(title, subtitle string, content fyne.CanvasObject) fyne.CanvasObject {
	// Создаем тень
	shadow := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 30})
	shadow.Resize(fyne.NewSize(0, 0))
	shadow.Move(fyne.NewPos(3, 3))

	// Создаем карточку
	card := widget.NewCard(title, subtitle, content)

	// Создаем контейнер с тенью
	return container.NewMax(shadow, card)
}

// GradientButton создает кнопку с градиентным фоном
type GradientButton struct {
	widget.BaseWidget
	text       string
	onTapped   func()
	startColor color.Color
	endColor   color.Color
}

func NewGradientButton(text string, onTapped func()) *GradientButton {
	btn := &GradientButton{
		text:       text,
		onTapped:   onTapped,
		startColor: ForestGreen,
		endColor:   DarkForestGreen,
	}
	btn.ExtendBaseWidget(btn)
	return btn
}

func (b *GradientButton) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		// Создаем вертикальный градиент
		ratio := float32(y) / float32(h)
		return lerpColor(b.startColor, b.endColor, ratio)
	})

	text := canvas.NewText(b.text, Cream)
	text.TextStyle = fyne.TextStyle{Bold: true}
	text.Alignment = fyne.TextAlignCenter

	return &gradientButtonRenderer{
		bg:     bg,
		text:   text,
		button: b,
	}
}

type gradientButtonRenderer struct {
	bg     *canvas.Raster
	text   *canvas.Text
	button *GradientButton
}

func (r *gradientButtonRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.text.Resize(size)
}

func (r *gradientButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(100, 40)
}

func (r *gradientButtonRenderer) Refresh() {
	r.text.Text = r.button.text
	r.text.Refresh()
	r.bg.Refresh()
}

func (r *gradientButtonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.text}
}

func (r *gradientButtonRenderer) Destroy() {}

func (b *GradientButton) Tapped(_ *fyne.PointEvent) {
	if b.onTapped != nil {
		b.onTapped()
	}
}

func (b *GradientButton) TappedSecondary(_ *fyne.PointEvent) {}

// IconButton создает кнопку с иконкой
func IconButton(icon fyne.Resource, text string, onTapped func()) *widget.Button {
	btn := widget.NewButtonWithIcon(text, icon, onTapped)
	return btn
}

// FloatingActionButton создает круглую кнопку действия
type FloatingActionButton struct {
	widget.BaseWidget
	icon     fyne.Resource
	onTapped func()
}

func NewFloatingActionButton(icon fyne.Resource, onTapped func()) *FloatingActionButton {
	fab := &FloatingActionButton{
		icon:     icon,
		onTapped: onTapped,
	}
	fab.ExtendBaseWidget(fab)
	return fab
}

func (f *FloatingActionButton) CreateRenderer() fyne.WidgetRenderer {
	// Создаем круглый фон
	bg := canvas.NewCircle(ForestGreen)

	// Создаем тень
	shadow := canvas.NewCircle(color.NRGBA{R: 0, G: 0, B: 0, A: 50})
	shadow.Move(fyne.NewPos(2, 2))

	// Создаем иконку
	icon := canvas.NewImageFromResource(f.icon)
	icon.FillMode = canvas.ImageFillContain

	return &fabRenderer{
		shadow: shadow,
		bg:     bg,
		icon:   icon,
		fab:    f,
	}
}

type fabRenderer struct {
	shadow *canvas.Circle
	bg     *canvas.Circle
	icon   *canvas.Image
	fab    *FloatingActionButton
}

func (r *fabRenderer) Layout(size fyne.Size) {
	// Делаем кнопку круглой
	diameter := fyne.Min(size.Width, size.Height)
	circleSize := fyne.NewSize(diameter, diameter)

	r.shadow.Resize(circleSize)
	r.bg.Resize(circleSize)

	// Размещаем иконку в центре с отступами
	iconSize := diameter * 0.6
	iconPos := (diameter - iconSize) / 2
	r.icon.Move(fyne.NewPos(iconPos, iconPos))
	r.icon.Resize(fyne.NewSize(iconSize, iconSize))
}

func (r *fabRenderer) MinSize() fyne.Size {
	return fyne.NewSize(56, 56)
}

func (r *fabRenderer) Refresh() {
	r.bg.FillColor = ForestGreen
	r.bg.Refresh()
}

func (r *fabRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.shadow, r.bg, r.icon}
}

func (r *fabRenderer) Destroy() {}

func (f *FloatingActionButton) Tapped(_ *fyne.PointEvent) {
	if f.onTapped != nil {
		f.onTapped()
	}
}

func (f *FloatingActionButton) TappedSecondary(_ *fyne.PointEvent) {}

// StyledEntry создает стилизованное поле ввода
func StyledEntry(placeholder string) *widget.Entry {
	entry := widget.NewEntry()
	entry.PlaceHolder = placeholder
	return entry
}

// InfoCard создает информационную карточку с иконкой
func InfoCard(icon fyne.Resource, title, value string, bgColor color.Color) fyne.CanvasObject {
	// Фон карточки
	bg := canvas.NewRectangle(bgColor)
	bg.CornerRadius = 8

	// Иконка
	iconImg := canvas.NewImageFromResource(icon)
	iconImg.FillMode = canvas.ImageFillContain
	iconImg.SetMinSize(fyne.NewSize(40, 40))

	// Тексты
	titleLabel := canvas.NewText(title, Cream)
	titleLabel.TextSize = 12

	valueLabel := canvas.NewText(value, Cream)
	valueLabel.TextSize = 24
	valueLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Компоновка
	content := container.NewBorder(
		nil, nil,
		container.NewCenter(iconImg),
		nil,
		container.NewVBox(
			titleLabel,
			valueLabel,
		),
	)

	// Добавляем отступы
	padded := container.NewPadded(content)

	return container.NewMax(bg, padded)
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
