package desktop

import (
	"time"

	"fyne.io/fyne/v2"

	"github.com/Lekuruu/go-puush-client/internal/config"
	"github.com/Lekuruu/go-puush-client/internal/hotkeys"
	"github.com/Lekuruu/go-puush-client/internal/tray"
	"github.com/Lekuruu/go-puush-client/pkg/puush"
)

// UI manages the desktop application windows and state.
type UI struct {
	app     fyne.App
	api     *puush.Client
	config  *config.Config
	tray    *tray.TrayManager
	hotkeys *hotkeys.HotkeyManager
}

func NewUI(app fyne.App, api *puush.Client, cfg *config.Config) *UI {
	tm := tray.NewTrayManager(api)
	hkm := hotkeys.NewHotkeyManager(cfg, tm)

	return &UI{
		app:     app,
		api:     api,
		config:  cfg,
		tray:    tm,
		hotkeys: hkm,
	}
}

func (ui *UI) Run() {
	// TODO: Maybe add some sort of theme customization?
	ui.app.Settings().SetTheme(NewWindowsTheme())

	// Show quickstart window if no credentials have been set
	// Otherwise, re-authenticate to see if the API key is still valid
	if !ui.api.Account.Credentials.HasApiKey() {
		ui.ShowStartupWindow()
	}

	// Initialize & start tray
	if ui.tray != nil {
		ui.tray.Initialize("puush")
		ui.tray.Apply(ui.app)
		ui.tray.SetSettingsCallback(ui.ShowSettingsWindow)
		ui.tray.SetPuushingDisabled(ui.config.General.DisabledToggle)

		// Apply configuration for copying to clipboard
		if ui.config.General.CopyToClipboard {
			ui.tray.EnableClipboard()
		} else {
			ui.tray.DisableClipboard()
		}

		// Save images locally, if enabled
		if ui.config.Capture.SaveImages {
			ui.tray.SetScreenshotsPath(ui.config.Capture.SaveImagePath)
		}

		// Start directory monitoring
		if len(ui.config.Capture.MonitorDirectories) > 0 {
			ui.tray.StartMonitor(ui.config.Capture.MonitorDirectories)
		}

		go ui.tray.PerformBackgroundAuthentication()
		go ui.tray.RefreshHistory()

		ui.hotkeys.Start()
	}

	ui.app.Run()
}

func (ui *UI) OnShutdown() {
	ui.tray.StopMonitor()

	// Keep "puushing disabled" state
	ui.config.General.DisabledToggle = ui.tray.PuushingDisabled()

	if ui.api.Account.Credentials.HasApiKey() {
		// Update account state in config after shutdown
		ui.config.Account.Key = *ui.api.Account.Credentials.Key
		ui.config.Account.Username = *ui.api.Account.Credentials.Identifier
		ui.config.Account.Type = int(ui.api.Account.Type)
		ui.config.Account.Usage = ui.api.Account.DiskUsage

		if ui.api.Account.SubscriptionEnd != nil {
			ui.config.Account.Expiry = ui.api.Account.SubscriptionEnd.Format(time.DateTime)
		}
	}
}
