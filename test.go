package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type BorderedButton struct {
	widget.BaseWidget
	Button *widget.Button
}

func NewBorderedButton(text string, tapped func()) *BorderedButton {
	b := &BorderedButton{
		Button: widget.NewButton(text, tapped),
	}
	b.ExtendBaseWidget(b)
	return b
}

func (b *BorderedButton) CreateRenderer() fyne.WidgetRenderer {
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeColor = color.NRGBA{R: 171, G: 173, B: 179, A: 255}
	border.StrokeWidth = 1
	
	c := container.NewStack(b.Button, border)
	return widget.NewSimpleRenderer(c)
}

func (b *BorderedButton) Disable() { b.Button.Disable() }
func (b *BorderedButton) Enable() { b.Button.Enable() }
func (b *BorderedButton) SetText(text string) { b.Button.SetText(text) }
func (b *BorderedButton) Text() string { return b.Button.Text }

func main() {
	a := app.New()
	w := a.NewWindow("Test")
	
	btn := NewBorderedButton("Click Me", func() {
		fmt.Println("Clicked!")
	})
	btn.Disable()
	
	w.SetContent(container.NewPadded(btn))
	
	w.Show()
}
