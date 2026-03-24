package tray

import (
	"io"
	"path/filepath"

	"github.com/Lekuruu/go-puush-client/assets"
	"github.com/Lekuruu/go-puush-client/internal/notifications"
	"github.com/Lekuruu/go-puush-client/pkg/puush"
)

func (m *TrayManager) PerformUpload(reader io.Reader, filename string) {
	url, err := m.api.Upload(reader, filename)
	if err != nil {
		m.OnUploadError(err)
		return
	}
	m.OnUploadComplete(url)
}

func (m *TrayManager) PerformProgressUpload(reader io.ReadCloser, total int64, filename string) {
	pr := puush.NewProgressReader(reader, total, m.OnTrayProgressUpdate)
	defer pr.Close()
	defer m.OnTrayProgressComplete()

	m.PerformUpload(pr, filename)
}

func (m *TrayManager) PerformFileUpload(path string) {
	pr, err := puush.NewProgressReaderFromFile(path, m.OnTrayProgressUpdate)
	if err != nil {
		m.OnUploadError(err)
		return
	}
	defer pr.Close()
	defer m.OnTrayProgressComplete()

	filename := filepath.Base(path)
	m.PerformUpload(pr, filename)
}

func (m *TrayManager) OnUploadComplete(url string) {
	// Update the tray icon to the "complete" state
	m.OnTrayProgressComplete()

	go notifications.NewNotification("puush complete!", "", url).
		WithSoundData(assets.SuccessSoundData).
		WithIconData(assets.PuushIconData).
		WithAction(url).
		Push()
	// TODO: Copy to clipboard depending on config
}

func (m *TrayManager) OnUploadError(err error) {
	// Update the tray icon to the "failed" state
	m.OnTrayProgressFail()

	go notifications.NewNotification("puush error", "", puush.FormatError(err)).
		WithIconData(assets.PuushIconData).
		Push()
	// TODO: Find right icon for error
}
