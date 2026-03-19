package desktop

import (
	"fyne.io/fyne/v2"
)

func (ui *UI) ShowSettingsWindow() {
	w := ui.app.NewWindow("Settings")
	w.Resize(fyne.NewSize(495, 360))
	w.Show()
}
