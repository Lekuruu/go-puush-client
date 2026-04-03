package tray

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.design/x/clipboard"
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

	filename := getImageFilename(reader)
	m.PerformUpload(reader, filename)
	m.OnScreenshotUploaded(reader, filename)
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

	filename := getImageFilename(reader)
	m.PerformUpload(reader, filename)
	m.OnScreenshotUploaded(reader, filename)
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

	filename := getImageFilename(reader)
	m.PerformUpload(reader, filename)
	m.OnScreenshotUploaded(reader, filename)
}

func (m *TrayManager) OnScreenshotUploaded(reader io.ReadSeeker, filename string) {
	if m.config.Capture.SaveImagesToClipboard {
		m.CopyScreenshotToClipboard(reader)
	}
	if m.config.Capture.SaveImages && m.config.Capture.SaveImagePath != "" {
		m.SaveScreenshotToDisk(reader, filename, m.config.Capture.SaveImagePath)
	}
}

func (m *TrayManager) CopyScreenshotToClipboard(reader io.ReadSeeker) {
	err := clipboard.Init()
	if err != nil {
		log.Printf("Error initializing clipboard: %v", err)
		return
	}

	reader.Seek(0, io.SeekStart)
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("Error reading screenshot data for clipboard: %v", err)
		return
	}

	clipboard.Write(clipboard.FmtImage, data)
	log.Printf("Screenshot image copied to clipboard")
}

func (m *TrayManager) SaveScreenshotToDisk(reader io.ReadSeeker, filename string, path string) string {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	// Seek back to start of reader, otherwise we aren't going to save anything
	reader.Seek(0, io.SeekStart)

	outputPath := path + filename
	outFile, err := os.Create(outputPath)
	if err != nil {
		log.Printf("Error creating file for saving screenshot: %v", err)
		return ""
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, reader)
	if err != nil {
		log.Printf("Error saving screenshot to file: %v", err)
		return ""
	}

	log.Printf("Screenshot saved to: %s", outputPath)
	return outputPath
}

func getImageFilename(reader io.ReadSeeker) string {
	ext := getImageExtension(reader)
	return fmt.Sprintf("ss (%s)%s", time.Now().Format("2006-01-02 at 15.04.05"), ext)
}

func getImageExtension(reader io.ReadSeeker) string {
	buffer := make([]byte, 512)
	n, _ := reader.Read(buffer)
	reader.Seek(0, io.SeekStart)

	contentType := http.DetectContentType(buffer[:n])
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ".png"
	}
}
