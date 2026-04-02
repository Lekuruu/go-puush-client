package tray

import (
	"fmt"
	"log"
	"time"
)

func (m *TrayManager) UploadAreaScreenshot() {
	provider := m.GetScreenshotProvider()
	if provider == nil {
		m.ShowErrorNotification("No screenshot provider available. Please install a compatible screenshot tool to use this feature!")
		return
	}

	reader, err := provider.CaptureArea()
	if err != nil {
		// TODO: Show notification to user if it's not a cancelled screeenshot
		log.Printf("Error capturing area screenshot: %v", err)
		return
	}
	defer reader.Close()

	filename := fmt.Sprintf("ss (%s).png", time.Now().Format("2006-01-02 at 15.04.05"))
	m.PerformUpload(reader, filename)
}

func (m *TrayManager) UploadDesktopScreenshot() {
	provider := m.GetScreenshotProvider()
	if provider == nil {
		m.ShowErrorNotification("No screenshot provider available. Please install a compatible screenshot tool to use this feature!")
		return
	}

	reader, err := provider.CaptureScreen()
	if err != nil {
		m.ShowErrorNotification("An error occurred while capturing the screenshot. Please try again.")
		return
	}
	defer reader.Close()

	filename := fmt.Sprintf("ss (%s).png", time.Now().Format("2006-01-02 at 15.04.05"))
	m.PerformUpload(reader, filename)
}

func (m *TrayManager) UploadWindowScreenshot() {
	provider := m.GetScreenshotProvider()
	if provider == nil {
		m.ShowErrorNotification("No screenshot provider available. Please install a compatible screenshot tool to use this feature!")
		return
	}

	reader, err := provider.CaptureWindow()
	if err != nil {
		m.ShowErrorNotification("An error occurred while capturing the screenshot. Please try again.")
		return
	}
	defer reader.Close()

	filename := fmt.Sprintf("ss (%s).png", time.Now().Format("2006-01-02 at 15.04.05"))
	m.PerformUpload(reader, filename)
}
