package tray

import (
	"github.com/sqweek/dialog"
)

func (m *TrayManager) UploadFileFromDialog() {
	filename, err := dialog.File().Title("Select a file to upload").Load()
	if err != nil {
		return
	}
	m.PerformFileUpload(filename)
}
