package desktop

import (
	"fyne.io/fyne/v2"
)

// UI manages the desktop application windows and state.
type UI struct {
	app fyne.App
}

func NewUI(app fyne.App) *UI {
	app.Settings().SetTheme(NewWindowsTheme())
	return &UI{app: app}
}

func (ui *UI) Run() {
	ui.ShowStartupWindow()
	ui.app.Run()
}
