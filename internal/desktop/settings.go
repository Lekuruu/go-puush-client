package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
)

func (ui *UI) ShowSettingsWindow() {
	w := ui.app.NewWindow("puush settings")
	w.Resize(fyne.NewSize(500, 350))
	w.SetIcon(puushIcon)

	tabs := container.NewAppTabs(
		container.NewTabItem("General", ui.buildGeneralTab()),
		container.NewTabItem("Key Bindings", ui.buildKeyBindingsTab()),
		container.NewTabItem("Account", widget.NewLabel("todo")),
		container.NewTabItem("Updates", widget.NewLabel("todo")),
		container.NewTabItem("Advanced", widget.NewLabel("todo")),
	)
	w.SetContent(container.NewPadded(tabs))
	w.Show()
}

func (ui *UI) buildGeneralTab() fyne.CanvasObject {
	startupCheckbox := widget.NewCheck("Start puush on startup", func(b bool) { ui.config.General.Startup = b })
	startupCheckbox.Checked = ui.config.General.Startup

	soundCheckbox := widget.NewCheck("Play a notification sound", func(b bool) { ui.config.General.NotificationSound = b })
	soundCheckbox.Checked = ui.config.General.NotificationSound

	copyLinkCheckbox := widget.NewCheck("Copy link to clipboard", func(b bool) { ui.config.General.CopyToClipboard = b })
	copyLinkCheckbox.Checked = ui.config.General.CopyToClipboard

	openBrowserCheckbox := widget.NewCheck("Open link in browser", func(b bool) { ui.config.General.OpenBrowser = b })
	openBrowserCheckbox.Checked = ui.config.General.OpenBrowser

	saveClipboardCheckbox := widget.NewCheck("Save image to the clipboard", func(b bool) { ui.config.Capture.SaveImagesToClipboard = b })
	saveClipboardCheckbox.Checked = ui.config.Capture.SaveImagesToClipboard

	saveLocalCheckbox := widget.NewCheck("Save a local copy of image", func(b bool) { ui.config.Capture.SaveImages = b })
	saveLocalCheckbox.Checked = ui.config.Capture.SaveImages

	savePathEntry := widget.NewEntry()
	savePathEntry.SetText(ui.config.Capture.SaveImagePath)
	savePathEntry.OnChanged = func(s string) { ui.config.Capture.SaveImagePath = s }

	browseButton := widget.NewButton("...", func() {
		path, err := dialog.Directory().Title("Select a file to upload").Browse()
		if err == nil {
			savePathEntry.SetText(path)
			ui.config.Capture.SaveImagePath = path
		}
	})
	saveLocalPathContainer := container.NewBorder(nil, nil, nil, browseButton, savePathEntry)

	onSuccessLeft := container.NewVBox(soundCheckbox, copyLinkCheckbox, openBrowserCheckbox)
	onSuccessRight := container.NewVBox(saveClipboardCheckbox, saveLocalCheckbox, saveLocalPathContainer)
	onSuccessGrid := container.NewGridWithColumns(2, onSuccessLeft, onSuccessRight)

	return container.NewVBox(
		widget.NewSeparator(),
		createGroup("General Settings", startupCheckbox),
		widget.NewSeparator(),
		createGroup("On successful puush", onSuccessGrid),
		widget.NewSeparator(),
	)
}

func (ui *UI) buildKeyBindingsTab() fyne.CanvasObject {
	fullScreenBtn := widget.NewButton("Ctrl+Shift+3", func() {})
	currentWindowBtn := widget.NewButton("Ctrl+Shift+2", func() {})
	captureAreaBtn := widget.NewButton("Ctrl+Shift+4", func() {})
	uploadFileBtn := widget.NewButton("Ctrl+Shift+U", func() {})
	uploadClipboardBtn := widget.NewButton("Ctrl+Shift+5", func() {})
	togglePuushBtn := widget.NewButton("Ctrl+Alt+P", func() {})

	rowFullscreen := container.NewGridWithColumns(2, widget.NewLabel("Capture full screen:"), fullScreenBtn)
	rowWindow := container.NewGridWithColumns(2, widget.NewLabel("Capture current window:"), currentWindowBtn)
	rowArea := container.NewGridWithColumns(2, widget.NewLabel("Capture Area:"), captureAreaBtn)
	rowFile := container.NewGridWithColumns(2, widget.NewLabel("Upload File:"), uploadFileBtn)
	// infoLabel := container.NewCenter(widget.NewLabel("Use this shortcut in Windows Explorer to quickly upload selected files."))
	rowClipboard := container.NewGridWithColumns(2, widget.NewLabel("Upload Clipboard:"), uploadClipboardBtn)
	rowToggle := container.NewGridWithColumns(2, widget.NewLabel("Toggle puush functionality:"), togglePuushBtn)

	content := container.NewVBox(
		rowFullscreen,
		rowWindow,
		rowArea,
		rowFile,
		rowClipboard,
		rowToggle,
	)

	return container.NewVBox(
		widget.NewSeparator(),
		createGroup("Keyboard Bindings", content),
		widget.NewSeparator(),
	)
}

func createGroup(title string, content fyne.CanvasObject) fyne.CanvasObject {
	indentedContent := container.NewBorder(nil, nil, widget.NewLabel("    "), widget.NewLabel("    "), content)
	return widget.NewCard("", title, indentedContent)
}
