// internal/ui/theme.go
package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ResortTheme - кастомная тема для базы отдыха
type ResortTheme struct{}

// Цветовая палитра
var (
	// Зеленые оттенки
	PrimaryGreen   = color.NRGBA{R: 45, G: 80, B: 22, A: 255}    // #2d5016
	SecondaryGreen = color.NRGBA{R: 74, G: 124, B: 35, A: 255}   // #4a7c23
	LightGreen     = color.NRGBA{R: 126, G: 160, B: 71, A: 255}  // #7ea047
	AccentGreen    = color.NRGBA{R: 156, G: 184, B: 111, A: 255} // #9cb86f
	PaleGreen      = color.NRGBA{R: 197, G: 212, B: 154, A: 255} // #c5d49a

	// Молочные тона
	Cream      = color.NRGBA{R: 250, G: 247, B: 240, A: 255} // #faf7f0
	LightCream = color.NRGBA{R: 253, G: 252, B: 247, A: 255} // #fdfcf7

	// Коричневые
	WarmBrown  = color.NRGBA{R: 139, G: 69, B: 19, A: 255}   // #8b4513
	LightBrown = color.NRGBA{R: 160, G: 82, B: 45, A: 255}   // #a0522d
	SoftBrown  = color.NRGBA{R: 222, G: 184, B: 135, A: 255} // #deb887

	// Текст и границы
	TextDark    = color.NRGBA{R: 44, G: 62, B: 45, A: 255}    // #2c3e2d
	TextLight   = color.NRGBA{R: 90, G: 107, B: 93, A: 255}   // #5a6b5d
	BorderLight = color.NRGBA{R: 232, G: 229, B: 212, A: 255} // #e8e5d4

	// Статусы
	SuccessColor = color.NRGBA{R: 74, G: 124, B: 35, A: 255}   // #4a7c23
	WarningColor = color.NRGBA{R: 212, G: 165, B: 116, A: 255} // #d4a574
	DangerColor  = color.NRGBA{R: 200, G: 90, B: 90, A: 255}   // #c85a5a
)

// Color возвращает цвет для конкретного использования
func (t ResortTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		return SecondaryGreen
	case theme.ColorNameBackground:
		return Cream
	case theme.ColorNameButton:
		return SecondaryGreen
	case theme.ColorNameHover:
		return AccentGreen
	case theme.ColorNamePressed:
		return PrimaryGreen
	case theme.ColorNameFocus:
		return LightGreen
	case theme.ColorNameSelection:
		return PaleGreen
	case theme.ColorNameForeground:
		return TextDark
	case theme.ColorNameDisabled:
		return TextLight
	case theme.ColorNamePlaceHolder:
		return TextLight
	case theme.ColorNameInputBackground:
		return color.White
	case theme.ColorNameInputBorder:
		return BorderLight
	case theme.ColorNameMenuBackground:
		return LightCream
	case theme.ColorNameOverlayBackground:
		return color.NRGBA{R: 45, G: 80, B: 22, A: 128} // Полупрозрачный зеленый
	case theme.ColorNameShadow:
		return color.NRGBA{R: 45, G: 80, B: 22, A: 25}
	case theme.ColorNameSuccess:
		return SuccessColor
	case theme.ColorNameWarning:
		return WarningColor
	case theme.ColorNameError:
		return DangerColor
	}

	// Возвращаем стандартную тему для неопределенных цветов
	return theme.DefaultTheme().Color(name, variant)
}

// Font возвращает шрифт для конкретного использования
func (t ResortTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon возвращает иконку
func (t ResortTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size возвращает размер для конкретного использования
func (t ResortTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNameScrollBar:
		return 16
	case theme.SizeNameSeparatorThickness:
		return 1
	case theme.SizeNameInputBorder:
		return 2
	}

	return theme.DefaultTheme().Size(name)
}

// Стили текста
type TextStyle int

const (
	TextStyleNormal TextStyle = iota
	TextStyleTitle
	TextStyleSubtitle
	TextStyleCaption
	TextStyleSuccess
	TextStyleWarning
	TextStyleError
)

// Кастомные функции для создания стилизованных виджетов
func CreateStyledCard(title, subtitle string, content fyne.CanvasObject) *fyne.Container {
	var cardContent []fyne.CanvasObject

	if title != "" {
		titleLabel := widget.NewLabel(title)
		titleLabel.TextStyle = fyne.TextStyle{Bold: true}
		cardContent = append(cardContent, titleLabel)
	}
	if subtitle != "" {
		subtitleLabel := widget.NewLabel(subtitle)
		cardContent = append(cardContent, subtitleLabel)
	}
	if content != nil {
		cardContent = append(cardContent, content)
	}

	card := container.NewVBox(cardContent...)

	// Создаем фон для карточки
	bg := canvas.NewRectangle(color.White)
	bg.CornerRadius = 12
	bg.StrokeWidth = 1
	bg.StrokeColor = BorderLight

	return container.NewStack(bg, container.NewPadded(card))
}

func NewStyledLabel(text string, style TextStyle) *widget.Label {
	label := widget.NewLabel(text)

	switch style {
	case TextStyleTitle:
		label.TextStyle = fyne.TextStyle{Bold: true}
	case TextStyleSubtitle:
		// Можно настроить дополнительные стили если нужно
	case TextStyleCaption:
		// Еще меньший размер
	case TextStyleSuccess:
		// Зеленый цвет для успеха
	case TextStyleWarning:
		// Желтый для предупреждений
	case TextStyleError:
		// Красный для ошибок
	}

	return label
}

// Кастомные кнопки
func NewPrimaryButton(text string, tapped func()) *widget.Button {
	btn := widget.NewButton(text, tapped)
	btn.Importance = widget.HighImportance
	return btn
}

func NewSecondaryButton(text string, tapped func()) *widget.Button {
	btn := widget.NewButton(text, tapped)
	btn.Importance = widget.MediumImportance
	return btn
}

// Статусные цвета для календаря
func GetStatusColor(status string) color.Color {
	switch status {
	case "free":
		return AccentGreen
	case "booked":
		return WarningColor
	case "occupied":
		return WarmBrown
	default:
		return BorderLight
	}
}

// Создание статусного индикатора
func CreateStatusIndicator(status string) *fyne.Container {
	statusColor := GetStatusColor(status)
	indicator := canvas.NewCircle(statusColor)
	indicator.Resize(fyne.NewSize(8, 8))

	return container.NewWithoutLayout(indicator)
}


