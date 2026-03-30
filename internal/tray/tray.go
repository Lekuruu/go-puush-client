package tray

import (
	"errors"
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/Lekuruu/go-puush-client/assets"
	"github.com/Lekuruu/go-puush-client/internal/config"
	"github.com/Lekuruu/go-puush-client/internal/notifications"
	"github.com/Lekuruu/go-puush-client/internal/screenshots"
	"github.com/Lekuruu/go-puush-client/pkg/puush"
)

type TrayManager struct {
	api         *puush.Client
	screenshots screenshots.ScreenshotProvider

	menu             *fyne.Menu
	targetApp        fyne.App
	settingsCallback func()
}

func NewTrayManager(api *puush.Client) *TrayManager {
	provider, _ := screenshots.GetDefaultProvider()
	return &TrayManager{api: api, screenshots: provider}
}

// SetSettingsCallback will set the function that will be called
// once the "Settings..." action has been invoked
func (m *TrayManager) SetSettingsCallback(callback func()) {
	m.settingsCallback = callback
}

// GetScreenshotProvider returns the screenshot provider used by the tray manager
func (m *TrayManager) GetScreenshotProvider() screenshots.ScreenshotProvider {
	return m.screenshots
}

// SetScreenshotProvider sets the screenshot provider for the tray manager
func (m *TrayManager) SetScreenshotProvider(provider screenshots.ScreenshotProvider) {
	m.screenshots = provider
}

// ShowUploadNotification will display a notification indicating that an upload was successful
func (m *TrayManager) ShowUploadNotification(url string) {
	go notifications.NewNotification("puush complete!", "", url).
		WithSoundData(assets.SuccessSoundData).
		WithIconData(assets.PuushIconData).
		WithAction(url).
		Push()
}

// ShowErrorNotification will display an error notification with the provided message
func (m *TrayManager) ShowErrorNotification(message string) {
	go notifications.NewNotification("puush error", "", message).
		WithIconData(assets.PuushIconData).
		Push()
	// TODO: Find right icon for error
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

	accountSettings := fyne.NewMenuItem("My Account", func() {
		if !m.api.Account.Credentials.HasApiKey() {
			return
		}

		path := fmt.Sprintf("/login/go/?k=%s", *m.api.Account.Credentials.Key)
		accountUrl, _ := url.Parse(m.api.FormatURL(path))
		fyne.CurrentApp().OpenURL(accountUrl)
	})

	// TODO: Implement recent uploads
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
	uploadClipboard.Icon = clipboardIcon

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
