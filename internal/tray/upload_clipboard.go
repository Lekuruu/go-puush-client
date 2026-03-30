package tray

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"fyne.io/fyne/v2"
)

func (m *TrayManager) UploadFromClipboard() {
	content := resolveClipboard()
	if content == "" {
		m.ShowErrorNotification("Your clipboard is empty or does not contain any text.")
		return
	}

	reader := strings.NewReader(content)
	filename := fmt.Sprintf("clipboard (%s).txt", time.Now().Format("2006-01-02 at 15.04.05"))

	m.PerformUpload(reader, filename)
}

func resolveClipboard() string {
	content := fyne.CurrentApp().Clipboard().Content()
	if content != "" {
		return content
	}

	// Fallback to wl-paste for linux/wayland systems
	cmd := exec.Command("wl-paste")
	out, err := cmd.Output()
	if err == nil && len(out) > 0 {
		return string(out)
	}

	// No clipboard for u :(
	return ""
}
