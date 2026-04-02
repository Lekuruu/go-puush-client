package tray

import (
	"io"
	"log"
	"net/url"
	"path/filepath"

	"fyne.io/fyne/v2"
	"github.com/Lekuruu/go-puush-client/pkg/puush"
)

func (m *TrayManager) PerformUpload(reader io.Reader, filename string) {
	if !m.api.Account.Credentials.HasApiKey() {
		return
	}
	log.Println("Starting upload:", filename)

	urlResponse, err := m.api.Upload(reader, filename)
	if err != nil {
		m.OnUploadError(err)
		return
	}
	m.OnUploadComplete(urlResponse)
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

func (m *TrayManager) OnUploadComplete(urlResponse string) {
	log.Println("Upload complete:", urlResponse)

	// Set updated disk usage to config
	m.config.Account.Usage = m.api.Account.DiskUsage

	// Update the tray icon to the "complete" state
	m.OnTrayProgressComplete()
	m.ShowUploadNotification(urlResponse)

	if m.config.General.CopyToClipboard {
		fyne.CurrentApp().Clipboard().SetContent(urlResponse)
	}
	if m.config.General.OpenBrowser {
		if u, err := url.Parse(urlResponse); err == nil {
			fyne.CurrentApp().OpenURL(u)
		}
	}

	// Refresh the history to reflect the new upload
	go m.RefreshHistory()
}

func (m *TrayManager) OnUploadError(err error) {
	log.Println("Upload error:", err)

	// Update the tray icon to the "failed" state
	m.OnTrayProgressFail()
	m.ShowErrorNotification(puush.FormatError(err))
}
