package desktop

import (
	"fyne.io/fyne/v2"

	"github.com/Lekuruu/go-puush-client/internal/config"
	"github.com/Lekuruu/go-puush-client/internal/tray"
	"github.com/Lekuruu/go-puush-client/pkg/puush"
)

// UI manages the desktop application windows and state.
type UI struct {
	app    fyne.App
	api    *puush.Client
	config *config.Config
	tray   *tray.TrayManager
}

func NewUI(app fyne.App, api *puush.Client, cfg *config.Config) *UI {
	return &UI{
		app:    app,
		api:    api,
		config: cfg,
		tray:   tray.NewTrayManager(api),
	}
}

func (ui *UI) Run() {
	// TODO: Maybe add some sort of theme customization?
	ui.app.Settings().SetTheme(NewWindowsTheme())

	// Initialize & start tray
	if ui.tray != nil {
		ui.tray.Initialize("puush")
		ui.tray.Apply(ui.app)
		ui.tray.SetSettingsCallback(ui.ShowSettingsWindow)
	}

	// Show quickstart window if no credentials have been set
	// Otherwise, re-authenticate to see if the API key is still valid
	if !ui.api.Account.Credentials.HasApiKey() {
		ui.ShowStartupWindow()
	} else {
		ui.tray.PerformBackgroundAuthentication()
	}

	// Apply configuration for copying to clipboard
	if ui.config.General.CopyToClipboard {
		ui.tray.EnableClipboard()
	} else {
		ui.tray.DisableClipboard()
	}

	ui.app.Run()
}
