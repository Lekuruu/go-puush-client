package tray

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/Lekuruu/go-puush-client/internal/config"
)

type TrayManager struct {
	app  desktop.App
	menu *fyne.Menu
}

func NewTrayManager(app fyne.App) *TrayManager {
	if desktopApp, ok := app.(desktop.App); ok {
		return &TrayManager{app: desktopApp}
	}
	return nil
}

func (m *TrayManager) Initialize(applicationName string) error {
	puushVersion := fyne.NewMenuItem(config.VersionString(), func() {})
	puushVersion.Disabled = true
	accountSettings := fyne.NewMenuItem("My Account", func() {})

	recentUploads := fyne.NewMenuItem("Recent Uploads", func() {})
	recentUploads.Disabled = true

	captureWindow := fyne.NewMenuItem("Capture Current Window", func() {})
	captureWindow.Icon = windowIcon
	captureDesktop := fyne.NewMenuItem("Capture Desktop", func() {})
	captureDesktop.Icon = fullscreenIcon
	captureArea := fyne.NewMenuItem("Capture Area", func() {})
	captureArea.Icon = selectionIcon
	uploadFile := fyne.NewMenuItem("Upload File", func() {})
	uploadFile.Icon = uploadIcon
	uploadClipboard := fyne.NewMenuItem("Upload Clipboard", func() {})

	var disablePuushing *fyne.MenuItem

	disablePuushing = fyne.NewMenuItem("Disable puushing", func() {
		disablePuushing.Checked = !disablePuushing.Checked
		m.menu.Refresh()
	})
	settings := fyne.NewMenuItem("Settings...", func() {})

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
	m.app.SetSystemTrayMenu(m.menu)
	m.app.SetSystemTrayIcon(puushIcon)
	return nil
}
