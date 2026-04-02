package config

import (
	"net/url"
	"time"
)

// Store defines the interface for loading and saving the
// configuration across different platforms.
type Store interface {
	Load() (*Config, error)
	Save(cfg *Config) error
}

// Config represents the application configuration.
type Config struct {
	Account AccountConfig
	General GeneralConfig
	Capture CaptureConfig
	Hotkeys HotkeyConfig
	Misc    MiscConfig
}

type AccountConfig struct {
	Username string
	Key      string
	Type     int
	Usage    int64
	Expiry   string
}

type GeneralConfig struct {
	OpenBrowser       bool
	NotificationSound bool
	CopyToClipboard   bool
	Startup           bool
	ContextMenu       bool
	DisabledToggle    bool
}

type CaptureConfig struct {
	UploadQuality      int
	FullscreenMode     int
	SelectionRectangle bool
	SaveImages         bool
	SaveImagePath      string
	MonitorDirectories []string
}

type HotkeyConfig struct {
	ScreenSelection         string
	FullscreenScreenshot    string
	CurrentWindowScreenshot string
	UploadFile              string
	UploadClipboard         string
	Toggle                  string
}

type MiscConfig struct {
	LastUpdate time.Time
	ServerURL  string
}

func (misc *MiscConfig) ParseServerURL() *url.URL {
	obj, err := url.Parse(misc.ServerURL)
	if err != nil {
		// Url seems to be invalid, revert back to default
		misc.ServerURL = "https://puush.me"
		obj, _ = url.Parse(misc.ServerURL)
	}

	// Ensure path is cleared
	obj.Path = ""
	return obj
}

// DefaultConfig returns a Config populated with default values.
func DefaultConfig() *Config {
	return &Config{
		General: GeneralConfig{
			OpenBrowser:       false,
			NotificationSound: true,
			CopyToClipboard:   true,
			Startup:           true,
			ContextMenu:       true,
			DisabledToggle:    false,
		},
		Capture: CaptureConfig{
			UploadQuality:      1,
			FullscreenMode:     0,
			SelectionRectangle: true,
			SaveImages:         false,
			SaveImagePath:      "",
			MonitorDirectories: []string{},
		},
		Hotkeys: HotkeyConfig{
			ScreenSelection:         "Ctrl+Shift+4",
			FullscreenScreenshot:    "Ctrl+Shift+3",
			CurrentWindowScreenshot: "Ctrl+Shift+2",
			UploadFile:              "Ctrl+Shift+U",
			UploadClipboard:         "Ctrl+Shift+5",
			Toggle:                  "Ctrl+Alt+P",
		},
		Misc: MiscConfig{
			ServerURL:  "https://puush.me",
			LastUpdate: time.Now(),
		},
	}
}
