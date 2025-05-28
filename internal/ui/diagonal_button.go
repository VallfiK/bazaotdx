package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// DiagonalButton представляет кнопку с диагональным разделением
type DiagonalButton struct {
	widget.BaseWidget

	leftColor   color.Color
	rightColor  color.Color
	leftTapped  func()
	rightTapped func()
	text        string

	// Кешируем объекты для рендеринга
	leftTriangle  *canvas.Rectangle
	rightTriangle *canvas.Rectangle
	textLabel     *canvas.Text
}

// NewDiagonalButton создает новую кнопку с диагональным разделением
func NewDiagonalButton(leftColor, rightColor color.Color, text string, leftTapped, rightTapped func()) *DiagonalButton {
	db := &DiagonalButton{
		leftColor:   leftColor,
		rightColor:  rightColor,
		leftTapped:  leftTapped,
		rightTapped: rightTapped,
		text:        text,
	}
	db.ExtendBaseWidget(db)

	// Создаем объекты сразу
	db.leftTriangle = canvas.NewRectangle(leftColor)
	db.rightTriangle = canvas.NewRectangle(rightColor)
	db.textLabel = canvas.NewText(text, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
	db.textLabel.TextSize = 8
	db.textLabel.Alignment = fyne.TextAlignCenter

	return db
}

func (db *DiagonalButton) Tapped(evt *fyne.PointEvent) {
	size := db.Size()
	if size.Width == 0 || size.Height == 0 {
		return
	}

	// Вычисляем, на какую часть кнопки кликнули
	// Диагональ идет из левого верхнего угла в правый нижний
	k := size.Height / size.Width
	expectedY := evt.Position.X * k

	if evt.Position.Y < expectedY {
		// Клик на левую часть (верхний треугольник)
		if db.leftTapped != nil {
			db.leftTapped()
		}
	} else {
		// Клик на правую часть (нижний треугольник)
		if db.rightTapped != nil {
			db.rightTapped()
		}
	}
}

// TappedSecondary обрабатывает правый клик
func (db *DiagonalButton) TappedSecondary(_ *fyne.PointEvent) {}

// CreateRenderer создает рендерер для кнопки
func (db *DiagonalButton) CreateRenderer() fyne.WidgetRenderer {
	return &diagonalButtonRenderer{
		base: db,
	}
}

type diagonalButtonRenderer struct {
	base *DiagonalButton
}

// Layout не нужен для этого типа рендерера
func (r *diagonalButtonRenderer) Layout(size fyne.Size) {
	if r.base.leftTriangle != nil {
		r.base.leftTriangle.Resize(size)
	}
	if r.base.rightTriangle != nil {
		r.base.rightTriangle.Resize(size)
	}
	if r.base.textLabel != nil {
		r.base.textLabel.Resize(size)
	}
}

// MinSize возвращает минимальный размер кнопки
func (r *diagonalButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(35, 60)
}

// Refresh обновляет рендерер
func (r *diagonalButtonRenderer) Refresh() {
	// Обновляем цвета
	if r.base.leftTriangle != nil {
		r.base.leftTriangle.FillColor = r.base.leftColor
		r.base.leftTriangle.Refresh()
	}
	if r.base.rightTriangle != nil {
		r.base.rightTriangle.FillColor = r.base.rightColor
		r.base.rightTriangle.Refresh()
	}
	if r.base.textLabel != nil {
		r.base.textLabel.Text = r.base.text
		r.base.textLabel.Refresh()
	}
}

// Objects возвращает объекты для рендеринга
func (r *diagonalButtonRenderer) Objects() []fyne.CanvasObject {
	size := r.base.Size()
	if size.Width <= 0 || size.Height <= 0 {
		size = r.MinSize()
	}

	// Создаем маски для треугольников с помощью canvas.Raster
	leftMask := canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		if w == 0 || h == 0 {
			return color.Transparent
		}

		// Вычисляем коэффициент наклона диагонали
		k := float32(h) / float32(w)
		expectedY := int(float32(x) * k)

		// Если точка выше диагонали - левый треугольник
		if y < expectedY {
			return r.base.leftColor
		}
		return color.Transparent
	})

	rightMask := canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		if w == 0 || h == 0 {
			return color.Transparent
		}

		k := float32(h) / float32(w)
		expectedY := int(float32(x) * k)

		// Если точка ниже или на диагонали - правый треугольник
		if y >= expectedY {
			return r.base.rightColor
		}
		return color.Transparent
	})

	// Устанавливаем размеры
	leftMask.Resize(size)
	rightMask.Resize(size)

	// Создаем текст с черным цветом для лучшей видимости
	textLabel := canvas.NewText(r.base.text, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
	textLabel.TextSize = 8
	textLabel.Alignment = fyne.TextAlignCenter
	textLabel.Resize(size)

	// Возвращаем объекты в правильном порядке
	return []fyne.CanvasObject{
		leftMask,
		rightMask,
		textLabel,
	}
}

// Destroy уничтожает рендерер
func (r *diagonalButtonRenderer) Destroy() {
	if r != nil {
		r.base = nil
	}
}
