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

func createGroup(title string, content fyne.CanvasObject) fyne.CanvasObject {
	indentedContent := container.NewBorder(nil, nil, widget.NewLabel("    "), widget.NewLabel("    "), content)
	return widget.NewCard("", title, indentedContent)
}
