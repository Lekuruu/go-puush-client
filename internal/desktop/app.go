package desktop

import (
	"fyne.io/fyne/v2"

	"github.com/Lekuruu/go-puush-client/internal/config"
)

// UI manages the desktop application windows and state.
type UI struct {
	app    fyne.App
	config *config.Config
}

func NewUI(app fyne.App, cfg *config.Config) *UI {
	app.Settings().SetTheme(NewWindowsTheme())
	return &UI{app: app, config: cfg}
}

func (ui *UI) Run() {
	ui.ShowStartupWindow()
	ui.app.Run()
}
