package desktop

import (
	"fyne.io/fyne/v2"

	"github.com/Lekuruu/go-puush-client/internal/config"
	"github.com/Lekuruu/go-puush-client/pkg/puush"
)

// UI manages the desktop application windows and state.
type UI struct {
	app    fyne.App
	api    *puush.Client
	config *config.Config
}

func NewUI(app fyne.App, api *puush.Client, cfg *config.Config) *UI {
	app.Settings().SetTheme(NewWindowsTheme())
	return &UI{app: app, api: api, config: cfg}
}

func (ui *UI) Run() {
	ui.ShowStartupWindow()
	ui.app.Run()
}
