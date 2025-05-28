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
	return db
}

func (db *DiagonalButton) Tapped(evt *fyne.PointEvent) {
	size := db.Size()

	// Вычисляем, на какую часть кнопки кликнули
	// Диагональ идет из левого верхнего угла в правый нижний
	// Если точка выше диагонали - левая часть, ниже - правая

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
func (r *diagonalButtonRenderer) Layout(size fyne.Size) {}

// MinSize возвращает минимальный размер кнопки
func (r *diagonalButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(35, 60)
}

// Refresh обновляет рендерер
func (r *diagonalButtonRenderer) Refresh() {
	canvas.Refresh(r.base)
}

// getTextColor возвращает подходящий цвет текста в зависимости от фона
func (r *diagonalButtonRenderer) getTextColor() color.Color {
	// Для лучшей читаемости используем черный текст
	// Можно также добавить логику выбора цвета в зависимости от яркости фона
	return color.NRGBA{R: 0, G: 0, B: 0, A: 255} // Черный текст
}

// Objects возвращает объекты для рендеринга
func (r *diagonalButtonRenderer) Objects() []fyne.CanvasObject {
	size := r.base.Size()
	if size.Width == 0 || size.Height == 0 {
		size = r.MinSize()
	}

	// Создаем левый треугольник
	leftTriangle := canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		// Вычисляем коэффициент наклона диагонали
		k := float32(h) / float32(w)
		expectedY := int(float32(x) * k)

		// Если точка выше диагонали - левый треугольник
		if y < expectedY {
			return r.base.leftColor
		}
		return color.Transparent
	})

	// Создаем правый треугольник
	rightTriangle := canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		k := float32(h) / float32(w)
		expectedY := int(float32(x) * k)

		// Если точка ниже диагонали - правый треугольник
		if y >= expectedY {
			return r.base.rightColor
		}
		return color.Transparent
	})

	// Создаем текст с контрастным цветом
	label := canvas.NewText(r.base.text, r.getTextColor())
	label.TextSize = 9
	label.Alignment = fyne.TextAlignCenter

	// Устанавливаем размеры
	leftTriangle.Resize(size)
	rightTriangle.Resize(size)
	label.Resize(size)

	return []fyne.CanvasObject{
		leftTriangle,
		rightTriangle,
		label,
	}
}

// Destroy уничтожает рендерер
func (r *diagonalButtonRenderer) Destroy() {
	if r != nil {
		r.base = nil
	}
}
