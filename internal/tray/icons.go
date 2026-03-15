package tray

import (
	"github.com/Lekuruu/go-puush-client/assets"

	"fyne.io/fyne/v2"
)

var (
	puushIcon      fyne.Resource = fyne.NewStaticResource("icon-puush.png", assets.PuushIconData)
	windowIcon     fyne.Resource = fyne.NewStaticResource("icon-window.png", assets.WindowIconData)
	fullscreenIcon fyne.Resource = fyne.NewStaticResource("icon-fullscreen.png", assets.FullscreenIconData)
	uploadIcon     fyne.Resource = fyne.NewStaticResource("icon-upload.png", assets.UploadIconData)
	selectionIcon  fyne.Resource = fyne.NewStaticResource("icon-selection.png", assets.SelectionIconData)
)
