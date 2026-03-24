package tray

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/Lekuruu/go-puush-client/internal/config"
	"github.com/Lekuruu/go-puush-client/pkg/puush"
)

type TrayManager struct {
	api              *puush.Client
	menu             *fyne.Menu
	targetApp        fyne.App
	settingsCallback func()
}

func NewTrayManager(api *puush.Client) *TrayManager {
	return &TrayManager{api: api}
}

// SetSettingsCallback will set the function that will be called
// once the "Settings..." action has been invoked
func (m *TrayManager) SetSettingsCallback(callback func()) {
	m.settingsCallback = callback
}

// Refresh will instruct the tray to update its menu.
func (m *TrayManager) Refresh() error {
	if m.menu == nil {
		return errors.New("tray was not initialized")
	}
	m.menu.Refresh()
	return nil
}

// Apply applies the tray menu to the specified app.
func (m *TrayManager) Apply(app fyne.App) error {
	if m.menu == nil {
		return errors.New("tray was not initialized")
	}
	if desktopApp, ok := app.(desktop.App); ok {
		desktopApp.SetSystemTrayMenu(m.menu)
		m.targetApp = app
		m.ResetTrayIcon()
		return nil
	}
	return errors.New("provided app is not a desktop app")
}

// Initialize populates the system tray menu.
func (m *TrayManager) Initialize(applicationName string) error {
	puushVersion := fyne.NewMenuItem(config.VersionString(), func() {})
	puushVersion.Disabled = true
	accountSettings := fyne.NewMenuItem("My Account", func() {})

	recentUploads := fyne.NewMenuItem("Recent Uploads", func() {})
	recentUploads.Disabled = true

	captureWindow := fyne.NewMenuItem("Capture Current Window", m.UploadWindowScreenshot)
	captureWindow.Icon = windowIcon
	captureDesktop := fyne.NewMenuItem("Capture Desktop", m.UploadDesktopScreenshot)
	captureDesktop.Icon = fullscreenIcon
	captureArea := fyne.NewMenuItem("Capture Area", m.UploadAreaScreenshot)
	captureArea.Icon = selectionIcon
	uploadFile := fyne.NewMenuItem("Upload File", m.UploadFileFromDialog)
	uploadFile.Icon = uploadIcon
	uploadClipboard := fyne.NewMenuItem("Upload Clipboard", m.UploadFromClipboard)

	var disablePuushing *fyne.MenuItem

	disablePuushing = fyne.NewMenuItem("Disable puushing", func() {
		disablePuushing.Checked = !disablePuushing.Checked
		m.menu.Refresh()
	})
	settings := fyne.NewMenuItem("Settings...", func() { m.settingsCallback() })

	m.menu = fyne.NewMenu(applicationName,
		puushVersion,
		accountSettings,
		fyne.NewMenuItemSeparator(),
		recentUploads,
		fyne.NewMenuItemSeparator(),
		captureWindow,
		captureDesktop,
		captureArea,
		uploadClipboard,
		uploadFile,
		fyne.NewMenuItemSeparator(),
		disablePuushing,
		settings,
	)
	return nil
}
