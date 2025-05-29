// internal/ui/theme.go
package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// ForestTheme - улучшенная кастомная тема "Звуки Леса"
type ForestTheme struct{}

// Обновленная цветовая палитра "Звуки Леса"
var (
	// Основные зеленые тона (более насыщенные и глубокие)
	ForestGreen      = color.NRGBA{R: 34, G: 139, B: 34, A: 255}   // Лесной зеленый
	DarkForestGreen  = color.NRGBA{R: 25, G: 100, B: 25, A: 255}   // Темно-зеленый
	LightForestGreen = color.NRGBA{R: 144, G: 238, B: 144, A: 255} // Светло-зеленый
	MintGreen        = color.NRGBA{R: 152, G: 251, B: 152, A: 255} // Мятный

	// Природные оттенки
	EarthBrown   = color.NRGBA{R: 139, G: 115, B: 85, A: 255}  // Земляной коричневый
	SkyBlue      = color.NRGBA{R: 135, G: 206, B: 235, A: 255} // Небесно-голубой
	SunsetOrange = color.NRGBA{R: 255, G: 160, B: 122, A: 255} // Закатный оранжевый

	// Белые и светлые тона
	White      = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // Чистый белый
	LightWhite = color.NRGBA{R: 250, G: 250, B: 250, A: 255} // Светло-белый
	CreamWhite = color.NRGBA{R: 253, G: 245, B: 230, A: 255} // Кремовый
	SoftWhite  = color.NRGBA{R: 248, G: 248, B: 255, A: 255} // Мягкий белый

	// Элементы интерфейса (более яркие и контрастные)
	ElementGreen   = color.NRGBA{R: 50, G: 205, B: 50, A: 255}   // Яркий зеленый
	ElementLight   = color.NRGBA{R: 173, G: 255, B: 173, A: 255} // Светлый элемент
	ElementHover   = color.NRGBA{R: 46, G: 125, B: 50, A: 255}   // Наведение
	ElementPressed = color.NRGBA{R: 34, G: 90, B: 34, A: 255}    // Нажатие

	// Акцентные цвета
	GoldenYellow = color.NRGBA{R: 255, G: 215, B: 0, A: 255}  // Золотой
	DeepRed      = color.NRGBA{R: 178, G: 34, B: 34, A: 255}  // Глубокий красный
	RoyalBlue    = color.NRGBA{R: 65, G: 105, B: 225, A: 255} // Королевский синий

	// Системные цвета
	SoftGray     = color.NRGBA{R: 245, G: 245, B: 245, A: 255} // Мягкий серый
	DarkGray     = color.NRGBA{R: 105, G: 105, B: 105, A: 255} // Темно-серый
	ErrorRed     = color.NRGBA{R: 220, G: 20, B: 60, A: 255}   // Красный для ошибок
	WarningAmber = color.NRGBA{R: 255, G: 191, B: 0, A: 255}   // Янтарный для предупреждений
	SuccessGreen = color.NRGBA{R: 0, G: 128, B: 0, A: 255}     // Зеленый для успеха
)

func (t *ForestTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	// Основные цвета интерфейса
	case theme.ColorNameBackground:
		return theme.DefaultTheme().Color(name, variant)

	case theme.ColorNameForeground:
		if variant == theme.VariantLight {
			return DarkForestGreen
		}
		return White

	case theme.ColorNameButton:
		if variant == theme.VariantLight {
			return ElementGreen // Более яркий зеленый для кнопок
		}
		return ForestGreen

	case theme.ColorNamePrimary:
		return ElementGreen

	case theme.ColorNameHover:
		if variant == theme.VariantLight {
			return ElementHover
		}
		return LightForestGreen

	case theme.ColorNamePressed:
		return ElementPressed

	case theme.ColorNameFocus:
		return GoldenYellow // Золотой для фокуса

	case theme.ColorNameSelection:
		return color.NRGBA{R: 173, G: 255, B: 173, A: 120} // Полупрозрачный светло-зеленый

	case theme.ColorNameInputBorder:
		if variant == theme.VariantLight {
			return ElementGreen
		}
		return LightForestGreen

	case theme.ColorNameShadow:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 80}

	case theme.ColorNameInputBackground:
		if variant == theme.VariantLight {
			return SoftWhite
		}
		return DarkForestGreen

	case theme.ColorNamePlaceHolder:
		if variant == theme.VariantLight {
			return color.NRGBA{R: 105, G: 105, B: 105, A: 180}
		}
		return color.NRGBA{R: 200, G: 200, B: 200, A: 180}

	case theme.ColorNameScrollBar:
		if variant == theme.VariantLight {
			return ElementLight
		}
		return ForestGreen

	case theme.ColorNameError:
		return ErrorRed

	case theme.ColorNameSuccess:
		return SuccessGreen

	case theme.ColorNameWarning:
		return WarningAmber

	case theme.ColorNameMenuBackground:
		if variant == theme.VariantLight {
			return SoftWhite
		}
		return DarkForestGreen

	case theme.ColorNameOverlayBackground:
		if variant == theme.VariantLight {
			return color.NRGBA{R: 248, G: 248, B: 255, A: 240}
		}
		return color.NRGBA{R: 25, G: 50, B: 25, A: 240}

	// Дополнительные цвета для карточек и панелей
	case theme.ColorNameHeaderBackground:
		return DarkForestGreen

	case theme.ColorNameSeparator:
		return ElementLight
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (t *ForestTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *ForestTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *ForestTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInlineIcon:
		return 24
	case theme.SizeNameScrollBar:
		return 16
	case theme.SizeNameScrollBarSmall:
		return 4
	case theme.SizeNameSeparatorThickness:
		return 2
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 26
	case theme.SizeNameSubHeadingText:
		return 18
	case theme.SizeNameCaptionText:
		return 12
	case theme.SizeNameInputBorder:
		return 2
	}
	return theme.DefaultTheme().Size(name)
}
