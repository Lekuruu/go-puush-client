package config

import (
	"net/url"
	"time"

	"github.com/Lekuruu/go-puush-client/internal/screenshots"
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
	UploadQuality      screenshots.Quality
	FullscreenMode     screenshots.FullscreenMode
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
			UploadQuality:      screenshots.QualityBest,
			FullscreenMode:     screenshots.FullscreenModeAllScreens,
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
