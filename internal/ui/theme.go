// internal/ui/theme.go
package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// ForestTheme - кастомная тема в стиле леса для базы отдыха
type ForestTheme struct{}

// Цветовая палитра
var (
	// Основные цвета
	ForestGreen      = color.NRGBA{R: 52, G: 140, B: 49, A: 255}   // Основной зеленый
	DarkForestGreen  = color.NRGBA{R: 29, G: 84, B: 28, A: 255}    // Темно-зеленый
	LightForestGreen = color.NRGBA{R: 129, G: 199, B: 127, A: 255} // Светло-зеленый
	MintGreen        = color.NRGBA{R: 183, G: 228, B: 182, A: 255} // Мятный

	// Белые оттенки
	White            = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // Белый
	LightWhite       = color.NRGBA{R: 250, G: 250, B: 250, A: 255} // Светло-белый
	SoftWhite        = color.NRGBA{R: 245, G: 245, B: 245, A: 255} // Мягкий белый

	// Зеленые оттенки для элементов
	ElementGreen     = color.NRGBA{R: 76, G: 175, B: 80, A: 255}   // Зеленый для элементов
	ElementLight     = color.NRGBA{R: 200, G: 230, B: 200, A: 255} // Светло-зеленый для элементов
	ElementHover     = color.NRGBA{R: 67, G: 160, B: 71, A: 255}   // Зеленый для наведения
	ElementPressed   = color.NRGBA{R: 55, G: 127, B: 60, A: 255}   // Зеленый для нажатия

	// Дополнительные цвета
	SoftGray     = color.NRGBA{R: 240, G: 240, B: 240, A: 255} // Мягкий серый
	DarkGray     = color.NRGBA{R: 100, G: 100, B: 100, A: 255} // Темно-серый
	ErrorRed     = color.NRGBA{R: 220, G: 53, B: 69, A: 255}   // Красный для ошибок
	WarningAmber = color.NRGBA{R: 255, G: 193, B: 7, A: 255}   // Янтарный для предупреждений
)

func (t *ForestTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	// Основные цвета интерфейса
	case theme.ColorNameBackground:
		if variant == theme.VariantLight {
			return White
		}
		return DarkForestGreen

	case theme.ColorNameForeground:
		if variant == theme.VariantLight {
			return DarkForestGreen
		}
		return White

	case theme.ColorNameButton:
		if variant == theme.VariantLight {
			return ElementGreen
		}
		return DarkForestGreen

	case theme.ColorNamePrimary:
		return ElementGreen

	case theme.ColorNameHover:
		if variant == theme.VariantLight {
			return ElementHover
		}
		return DarkForestGreen

	case theme.ColorNamePressed:
		return ElementPressed

	case theme.ColorNameFocus:
		return ElementGreen

	case theme.ColorNameSelection:
		return color.NRGBA{R: 129, G: 199, B: 127, A: 100} // Полупрозрачный зеленый

	case theme.ColorNameInputBorder:
		if variant == theme.VariantLight {
			return ElementLight
		}
		return DarkForestGreen

	case theme.ColorNameShadow:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 66}

	case theme.ColorNameInputBackground:
		if variant == theme.VariantLight {
			return LightWhite
		}
		return DarkForestGreen

	case theme.ColorNamePlaceHolder:
		if variant == theme.VariantLight {
			return color.NRGBA{R: 100, G: 100, B: 100, A: 180}
		}
		return color.NRGBA{R: 200, G: 200, B: 200, A: 180}

	case theme.ColorNameScrollBar:
		if variant == theme.VariantLight {
			return ElementLight
		}
		return DarkForestGreen

	case theme.ColorNameError:
		return ErrorRed

	case theme.ColorNameSuccess:
		return ForestGreen

	case theme.ColorNameWarning:
		return WarningAmber

	case theme.ColorNameMenuBackground:
		if variant == theme.VariantLight {
			return LightWhite
		}
		return DarkForestGreen

	case theme.ColorNameOverlayBackground:
		if variant == theme.VariantLight {
			return color.NRGBA{R: 250, G: 250, B: 250, A: 240}
		}
		return color.NRGBA{R: 40, G: 30, B: 20, A: 240}
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (t *ForestTheme) Font(style fyne.TextStyle) fyne.Resource {
	// Используем стандартные шрифты, но можно подключить кастомные
	return theme.DefaultTheme().Font(style)
}

func (t *ForestTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	// Используем стандартные иконки
	return theme.DefaultTheme().Icon(name)
}

func (t *ForestTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 6
	case theme.SizeNameInlineIcon:
		return 24
	case theme.SizeNameScrollBar:
		return 16
	case theme.SizeNameScrollBarSmall:
		return 3
	case theme.SizeNameSeparatorThickness:
		return 1
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 24
	case theme.SizeNameSubHeadingText:
		return 18
	case theme.SizeNameCaptionText:
		return 11
	case theme.SizeNameInputBorder:
		return 2
	}
	return theme.DefaultTheme().Size(name)
}
