package ui

import (
	"image/color"
	"log"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// DiagonalButton представляет кнопку с диагональным разделением
type DiagonalButton struct {
	widget.BaseWidget

	leftImage   string
	rightImage  string
	leftTapped  func()
	rightTapped func()
	text        string
	minSize     fyne.Size
	maxSize     fyne.Size

	// Кешируем объекты для рендеринга
	leftImageObj  *canvas.Image
	rightImageObj *canvas.Image
	textLabel     *canvas.Text
	bg            *canvas.Rectangle
}

// SetMinSize устанавливает минимальный размер кнопки
func (db *DiagonalButton) SetMinSize(size fyne.Size) {
	db.minSize = size
}

// SetMaxSize устанавливает максимальный размер кнопки
func (db *DiagonalButton) SetMaxSize(size fyne.Size) {
	db.maxSize = size
}

// MinSize возвращает минимальный размер кнопки
func (db *DiagonalButton) MinSize() fyne.Size {
	if db.minSize.Width > 0 && db.minSize.Height > 0 {
		return db.minSize
	}
	return fyne.NewSize(191, 62)
}

// NewDiagonalButton создает новую кнопку с диагональным разделением
func NewDiagonalButtonImage(leftImage, rightImage string, text string, leftTapped, rightTapped func()) *DiagonalButton {
	// Получаем путь к корневой директории проекта
	projectDir := "C:\\Users\\VallfIK\\Documents\\GitHub\\bazaotdx"
	leftImagePath := filepath.Join(projectDir, "images", leftImage)
	rightImagePath := filepath.Join(projectDir, "images", rightImage)

	// Логируем пути к изображениям
	log.Printf("📄 Путь к левому изображению: %s", leftImagePath)
	log.Printf("📄 Путь к правому изображению: %s", rightImagePath)

	// Проверяем существование файлов
	if _, err := os.Stat(leftImagePath); err != nil {
		if os.IsNotExist(err) {
			log.Printf("❌ Ошибка: Файл не найден: %s", leftImagePath)
			return nil
		} else {
			log.Printf("❌ Ошибка при проверке файла: %v", err)
			return nil
		}
	}

	if _, err := os.Stat(rightImagePath); err != nil {
		if os.IsNotExist(err) {
			log.Printf("❌ Ошибка: Файл не найден: %s", rightImagePath)
			return nil
		} else {
			log.Printf("❌ Ошибка при проверке файла: %v", err)
			return nil
		}
	}

	// Проверяем формат файлов
	if !strings.HasSuffix(leftImagePath, ".png") {
		log.Printf("❌ Ошибка: Файл %s не имеет расширения .png", leftImagePath)
		return nil
	}
	if !strings.HasSuffix(rightImagePath, ".png") {
		log.Printf("❌ Ошибка: Файл %s не имеет расширения .png", rightImagePath)
		return nil
	}

	log.Printf("✅ Проверка файлов успешно пройдена")

	db := &DiagonalButton{
		leftImage:   leftImagePath,
		rightImage:  rightImagePath,
		leftTapped:  leftTapped,
		rightTapped: rightTapped,
		text:        text,
	}
	db.ExtendBaseWidget(db)

	// Создаем объекты сразу
	db.leftImageObj = canvas.NewImageFromFile(leftImagePath)
	if db.leftImageObj == nil {
		log.Printf("❌ Ошибка: Не удалось загрузить изображение: %s", leftImagePath)
		return nil
	}
	log.Printf("✅ Изображение успешно загружено: %s", leftImagePath)

	db.rightImageObj = canvas.NewImageFromFile(rightImagePath)
	if db.rightImageObj == nil {
		log.Printf("❌ Ошибка: Не удалось загрузить изображение: %s", rightImagePath)
		return nil
	}
	log.Printf("✅ Изображение успешно загружено: %s", rightImagePath)

	// Устанавливаем размеры изображений
	if db.leftImageObj != nil {
		originalSize := db.leftImageObj.Size()
		log.Printf("✅ Размер оригинала левого изображения: %v", originalSize)
		db.leftImageObj.Resize(fyne.NewSize(191, 62))
		log.Printf("✅ Изображение успешно изменено размером")
	}

	if db.rightImageObj != nil {
		originalSize := db.rightImageObj.Size()
		log.Printf("✅ Размер оригинала правого изображения: %v", originalSize)
		db.rightImageObj.Resize(fyne.NewSize(191, 62))
		log.Printf("✅ Изображение успешно изменено размером")
	}

	// Проверяем размеры после изменения
	if db.leftImageObj != nil {
		newSize := db.leftImageObj.Size()
		log.Printf("✅ Новый размер левого изображения: %v", newSize)
	}

	if db.rightImageObj != nil {
		newSize := db.rightImageObj.Size()
		log.Printf("✅ Новый размер правого изображения: %v", newSize)
	}

	// Создаем белый фон для кнопки
	bgColor := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	db.bg = canvas.NewRectangle(bgColor)
	db.bg.Resize(fyne.NewSize(191, 62))

	// Создаем текст с более крупным размером и зеленым цветом
	textColor := color.NRGBA{R: 0, G: 128, B: 0, A: 255} // Зеленый цвет
	db.textLabel = canvas.NewText(text, textColor)
	db.textLabel.TextSize = 12 // Увеличиваем размер текста
	db.textLabel.Alignment = fyne.TextAlignCenter
	db.textLabel.Resize(fyne.NewSize(191, 62))

	log.Printf("✅ Кнопка успешно создана")
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
	bg   *canvas.Rectangle
}

// Layout не нужен для этого типа рендерера
func (r *diagonalButtonRenderer) Layout(size fyne.Size) {
	if r.bg != nil {
		r.bg.Resize(size)
	}
}

// MinSize возвращает минимальный размер кнопки
func (r *diagonalButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(191, 62)
}

// Refresh обновляет рендерер
func (r *diagonalButtonRenderer) Refresh() {
	if r.bg != nil {
		r.bg.FillColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // Белый фон
	}
	if r.base.leftImageObj != nil {
		r.base.leftImageObj.Refresh()
	}
	if r.base.rightImageObj != nil {
		r.base.rightImageObj.Refresh()
	}
	if r.base.textLabel != nil {
		r.base.textLabel.Refresh()
	}
}

// Objects возвращает объекты для рендеринга
func (r *diagonalButtonRenderer) Objects() []fyne.CanvasObject {
	size := r.base.Size()
	if size.Width <= 0 || size.Height <= 0 {
		size = r.MinSize()
	}

	// Создаем фон кнопки
	if r.bg == nil {
		r.bg = canvas.NewRectangle(color.White)
	}
	
	// Устанавливаем размеры для всех объектов
	if r.bg != nil {
		r.bg.Resize(size)
	}

	if r.base.leftImageObj != nil {
		r.base.leftImageObj.Resize(size)
	}

	if r.base.rightImageObj != nil {
		r.base.rightImageObj.Resize(size)
	}

	if r.base.textLabel != nil {
		r.base.textLabel.Resize(size)
		r.base.textLabel.Alignment = fyne.TextAlignCenter
		r.base.textLabel.TextSize = 12 // Увеличиваем размер текста
	}

	// Создаем маску для диагонали
	diagonalMask := canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		if w == 0 || h == 0 {
			return color.Transparent
		}

		// Вычисляем коэффициент наклона диагонали
		k := float32(h) / float32(w)
		expectedY := int(float32(x) * k)

		// Если точка выше диагонали - верхняя часть
		if y < expectedY {
			return color.White
		}
		// Если точка ниже или на диагонали - нижняя часть
		return color.Transparent
	})

	// Возвращаем все объекты в правильном порядке
	return []fyne.CanvasObject{
		r.bg,                  // Фон
		diagonalMask,         // Маска диагонали
		r.base.leftImageObj,  // Левое изображение
		r.base.rightImageObj, // Правое изображение
		r.base.textLabel,     // Текст
	}
}

// Destroy уничтожает рендерер
func (r *diagonalButtonRenderer) Destroy() {
	// В Fyne v2 не нужно вручную уничтожать объекты, они будут автоматически очищены сборщиком мусора
}
