package tray

import (
	"io"
	"path/filepath"

	"fyne.io/fyne/v2"
	"github.com/Lekuruu/go-puush-client/pkg/puush"
)

func (m *TrayManager) PerformUpload(reader io.Reader, filename string) {
	if !m.api.Account.Credentials.HasApiKey() {
		// TODO: Open startup window?
		return
	}

	url, err := m.api.Upload(reader, filename)
	if err != nil {
		m.OnUploadError(err)
		return
	}
	m.OnUploadComplete(url)
	// TODO: Implement upload retries
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
	m.ShowUploadNotification(url)

	if m.config.General.CopyToClipboard {
		fyne.CurrentApp().Clipboard().SetContent(url)
	}
}

func (m *TrayManager) OnUploadError(err error) {
	// Update the tray icon to the "failed" state
	m.OnTrayProgressFail()
	m.ShowErrorNotification(puush.FormatError(err))
}
