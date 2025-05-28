package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// DiagonalButton представляет кнопку с диагональным разделением
// Кнопка может обрабатывать клики на разных сторонах
// Если клик на левой стороне - вызывается leftTapped
// Если клик на правой стороне - вызывается rightTapped
type DiagonalButton struct {
	widget.BaseWidget

	leftColor   color.Color
	rightColor  color.Color
	leftTapped  func()
	rightTapped func()
	label       *canvas.Text
}

// NewDiagonalButton создает новую кнопку с диагональным разделением
func NewDiagonalButton(leftColor, rightColor color.Color, text string, leftTapped, rightTapped func()) *DiagonalButton {
	db := &DiagonalButton{
		leftColor:   leftColor,
		rightColor:  rightColor,
		leftTapped:  leftTapped,
		rightTapped: rightTapped,
		label:      canvas.NewText(text, color.White),
	}
	db.label.TextSize = 9
	db.label.Alignment = fyne.TextAlignCenter
	db.ExtendBaseWidget(db)
	return db
}

// CreateRenderer создает рендерер для кнопки
func (db *DiagonalButton) CreateRenderer() fyne.WidgetRenderer {
	// Создаем два прямоугольника для левой и правой стороны
	leftRect := canvas.NewRectangle(db.leftColor)
	rightRect := canvas.NewRectangle(db.rightColor)
	
	// Создаем линию для разделения
	line := canvas.NewLine(color.White)
	line.StrokeWidth = 2

	return &diagonalButtonRenderer{
		base:         db,
		leftRect:     leftRect,
		rightRect:    rightRect,
		line:         line,
		label:        db.label,
	}
}

type diagonalButtonRenderer struct {
	base         *DiagonalButton
	leftRect     *canvas.Rectangle
	rightRect    *canvas.Rectangle
	line         *canvas.Line
	label        *canvas.Text
}

// Layout устанавливает размеры и позиции элементов
func (r *diagonalButtonRenderer) Layout(size fyne.Size) {
	// Устанавливаем размеры прямоугольников
	width := size.Width
	height := size.Height

	// Левый прямоугольник занимает всю ширину
	r.leftRect.Resize(size)

	// Правый прямоугольник занимает всю ширину
	r.rightRect.Resize(size)

	// Линия по диагонали
	r.line.Position1 = fyne.NewPos(0, 0)
	r.line.Position2 = fyne.NewPos(width, height)

	// Устанавливаем размер текста
	r.label.Resize(size)
}

// MinSize возвращает минимальный размер кнопки
func (r *diagonalButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(35, 60) // Минимальный размер кнопки
}



// Destroy уничтожает рендерер
func (r *diagonalButtonRenderer) Destroy() {
	if r.leftRect != nil {
		r.leftRect = nil
	}
	if r.rightRect != nil {
		r.rightRect = nil
	}
	if r.line != nil {
		r.line = nil
	}
	if r.label != nil {
		r.label = nil
	}
}

// Refresh обновляет рендерер
func (r *diagonalButtonRenderer) Refresh() {
	canvas.Refresh(r.base)
}

// Objects возвращает объекты для рендеринга
func (r *diagonalButtonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		r.leftRect,
		r.rightRect,
		r.line,
		r.label,
	}
}

// Tapped обрабатывает клик по кнопке
func (db *DiagonalButton) Tapped(e *fyne.PointEvent) {
	// Проверяем, на какой стороне был клик
	width := db.Size().Width
	height := db.Size().Height
	
	// Вычисляем коэффициент наклона диагонали
	k := height / width
	
	// Вычисляем y для x клика по диагонали
	expectedY := k * e.AbsolutePosition.X
	
	// Если клик выше диагонали - левая сторона, иначе - правая
	if e.AbsolutePosition.Y < expectedY {
		if db.leftTapped != nil {
			db.leftTapped()
		}
	} else {
		if db.rightTapped != nil {
			db.rightTapped()
		}
	}
}

// TappedSecondary обрабатывает правый клик
func (db *DiagonalButton) TappedSecondary(e *fyne.PointEvent) {
	// Можно добавить обработку правого клика, если нужно
}
