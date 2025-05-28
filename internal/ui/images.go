package ui

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
)

func createImage(color color.Color, text string, filename string) error {
	// Создаем новый изображение 32x32 пикселя
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	
	// Заполняем фон
	draw.Draw(img, img.Bounds(), &image.Uniform{color}, image.Point{}, draw.Src)
	
	// Создаем файл
	file, err := os.Create(filepath.Join("..", "images", filename))
	if err != nil {
		return err
	}
	defer file.Close()
	
	// Сохраняем изображение
	return png.Encode(file, img)
}

func init() {
	// Создаем изображения с разными цветами
	createImage(color.RGBA{R: 0, G: 255, B: 0, A: 255}, "Св", "free.png")
	createImage(color.RGBA{R: 255, G: 0, B: 0, A: 255}, "З", "booked.png")
	createImage(color.RGBA{R: 0, G: 255, B: 0, A: 255}, "Св", "freefirst.png")
}
