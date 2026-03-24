package tray

import (
	"github.com/Lekuruu/go-puush-client/assets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

var (
	puushIcon      fyne.Resource = fyne.NewStaticResource("icon-puush.png", assets.PuushIconData)
	windowIcon     fyne.Resource = fyne.NewStaticResource("icon-window.png", assets.WindowIconData)
	fullscreenIcon fyne.Resource = fyne.NewStaticResource("icon-fullscreen.png", assets.FullscreenIconData)
	uploadIcon     fyne.Resource = fyne.NewStaticResource("icon-upload.png", assets.UploadIconData)
	selectionIcon  fyne.Resource = fyne.NewStaticResource("icon-selection.png", assets.SelectionIconData)
)

// OnTrayProgressUpdate gets called once the upload percentage
// has been updated through `puush.ProgressReader`
func (m *TrayManager) OnTrayProgressUpdate(percentage float64) {
	// TODO: ...
}

// OnTrayProgressComplete will reset the puush tray icon
// back to its original state
func (m *TrayManager) OnTrayProgressComplete() {
	if m.targetApp != nil {
		m.targetApp.(desktop.App).SetSystemTrayIcon(puushIcon)
	}
}
