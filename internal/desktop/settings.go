package desktop

import (
	"fyne.io/fyne/v2"
)

func (u *UI) ShowSettingsWindow() {
	w := u.app.NewWindow("Settings")
	w.Resize(fyne.NewSize(495, 360))
	w.Show()
}
