package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (ui *UI) ShowSettingsWindow() {
	w := ui.app.NewWindow("puush settings")
	w.Resize(fyne.NewSize(500, 350))
	w.SetIcon(puushIcon)

	generalView := ui.buildGeneralTab()
	keyBindingsView := ui.buildKeyBindingsTab()
	accountView, accountViewUpdate := ui.buildAccountTab()
	advancedView := ui.buildAdvancedTab(accountViewUpdate)

	tabs := container.NewAppTabs(
		container.NewTabItem("General", generalView),
		container.NewTabItem("Key Bindings", keyBindingsView),
		container.NewTabItem("Account", accountView),
		container.NewTabItem("Advanced", advancedView),
	)
	w.SetContent(container.NewPadded(tabs))
	w.Show()
}

func createGroup(title string, content fyne.CanvasObject) fyne.CanvasObject {
	indentedContent := container.NewBorder(nil, nil, widget.NewLabel("    "), widget.NewLabel("    "), content)
	return widget.NewCard("", title, indentedContent)
}

func createGroupNoIndent(title string, content fyne.CanvasObject) fyne.CanvasObject {
	return widget.NewCard("", title, content)
}

func trailingLabel(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.Alignment = fyne.TextAlignTrailing
	return label
}
