package desktop

import (
	"fyne.io/fyne/v2"
)

// UI manages the desktop application windows and state.
type UI struct {
	app fyne.App
}

func NewUI(a fyne.App) *UI {
	return &UI{app: a}
}

func (u *UI) Run() {
	u.ShowStartupWindow()
	u.app.Run()
}
