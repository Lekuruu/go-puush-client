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
	createHotkeyButton := func(initial string, onChange func(string)) *HotkeyButton {
		btn := NewHotkeyButton(initial)
		btn.OnStart = func() {
			ui.hotkeys.Stop()
			if canvas := ui.app.Driver().CanvasForObject(btn); canvas != nil {
				canvas.Focus(btn)
			}
		}
		btn.OnCancelled = func() {
			ui.hotkeys.Start()
			if canvas := ui.app.Driver().CanvasForObject(btn); canvas != nil {
				canvas.Focus(nil)
			}
		}
		btn.OnChanged = func(s string) {
			onChange(s)
			ui.hotkeys.Start()
			if canvas := ui.app.Driver().CanvasForObject(btn); canvas != nil {
				canvas.Focus(nil)
			}
		}
		return btn
	}

	fullScreenButton := createHotkeyButton(ui.config.Hotkeys.FullscreenScreenshot, func(s string) {
		ui.config.Hotkeys.FullscreenScreenshot = s
	})
	currentWindowButton := createHotkeyButton(ui.config.Hotkeys.CurrentWindowScreenshot, func(s string) {
		ui.config.Hotkeys.CurrentWindowScreenshot = s
	})
	captureAreaButton := createHotkeyButton(ui.config.Hotkeys.ScreenSelection, func(s string) {
		ui.config.Hotkeys.ScreenSelection = s
	})
	uploadFileButton := createHotkeyButton(ui.config.Hotkeys.UploadFile, func(s string) {
		ui.config.Hotkeys.UploadFile = s
	})
	uploadClipboardButton := createHotkeyButton(ui.config.Hotkeys.UploadClipboard, func(s string) {
		ui.config.Hotkeys.UploadClipboard = s
	})
	togglePuushButton := createHotkeyButton(ui.config.Hotkeys.Toggle, func(s string) {
		ui.config.Hotkeys.Toggle = s
	})

	rowFullscreen := container.NewGridWithColumns(2, widget.NewLabel("Capture full screen:"), fullScreenButton)
	rowWindow := container.NewGridWithColumns(2, widget.NewLabel("Capture current window:"), currentWindowButton)
	rowArea := container.NewGridWithColumns(2, widget.NewLabel("Capture Area:"), captureAreaButton)
	rowFile := container.NewGridWithColumns(2, widget.NewLabel("Upload File:"), uploadFileButton)
	rowClipboard := container.NewGridWithColumns(2, widget.NewLabel("Upload Clipboard:"), uploadClipboardButton)
	rowToggle := container.NewGridWithColumns(2, widget.NewLabel("Toggle puush functionality:"), togglePuushButton)

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
