package tray

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
)

func (m *TrayManager) UploadFromClipboard() {
	// NOTE: This should work in theory, but there's probably some wayland shenanigans that doens't allow it idk
	content := fyne.CurrentApp().Clipboard().Content()

	if content == "" {
		m.ShowErrorNotification("Your clipboard is empty or does not contain any text.")
		return
	}

	reader := strings.NewReader(content)
	filename := fmt.Sprintf("clipboard (%s).txt", time.Now().Format("2006-01-02 at 15.04.05"))

	m.PerformUpload(reader, filename)
}
