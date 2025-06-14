package background

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func NewCardBackground(
	backgroundColor color.Color,
	size fyne.Size,
) *canvas.Rectangle {
	bg := canvas.NewRectangle(backgroundColor)
	bg.SetMinSize(size)
	return bg
}

func NewBackgroundImage(imagePath string) *canvas.Image {
	bg := canvas.NewImageFromFile(imagePath)
	bg.FillMode = canvas.ImageFillStretch
	return bg
}