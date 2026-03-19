package desktop

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// newBorderedButton creates a regular button with a 1px border
func newBorderedButton(text string, callback func()) *fyne.Container {
	// Create a transparent rectangle with a solid border
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeColor = color.NRGBA{R: 171, G: 173, B: 179, A: 255}
	border.StrokeWidth = 1

	// Create button & place it on top of it
	button := widget.NewButton(text, callback)
	return container.NewStack(button, border)
}
