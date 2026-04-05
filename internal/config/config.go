package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/Lekuruu/go-puush-client/internal/screenshots"
	"github.com/Lekuruu/go-puush-client/pkg/puush"
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
	Type     puush.AccountType
	Usage    int64
	Expiry   string
}

func (a *AccountConfig) DiskUsageHumanReadable() string {
	return formatBytes(a.Usage)
}

func (a *AccountConfig) HasCredentials() bool {
	return a.Username != "" && a.Key != ""
}

func (a *AccountConfig) Reset() {
	a.Username = ""
	a.Key = ""
	a.Type = puush.AccountTypeRegular
	a.Usage = 0
	a.Expiry = ""
}

func (a *AccountConfig) SubscriptionExpiry() *time.Time {
	if a.Expiry == "" {
		return nil
	}
	expiryTime, err := time.Parse(time.RFC3339, a.Expiry)
	if err != nil {
		return nil
	}
	if expiryTime.IsZero() {
		return nil
	}
	return &expiryTime
}

type GeneralConfig struct {
	OpenBrowser       bool
	NotificationSound bool
	Startup           bool
	ContextMenu       bool
	DisabledToggle    bool
	CopyToClipboard   bool
}

type CaptureConfig struct {
	UploadQuality         screenshots.Quality
	FullscreenMode        screenshots.FullscreenMode
	SaveImages            bool
	SaveImagesToClipboard bool
	SaveImagePath         string
	MonitorDirectories    []string
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
			UploadQuality:         screenshots.QualityBest,
			FullscreenMode:        screenshots.FullscreenModeAllScreens,
			SaveImages:            false,
			SaveImagesToClipboard: false,
			SaveImagePath:         "",
			MonitorDirectories:    []string{},
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

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB", "PB", "EB"}
	return fmt.Sprintf("%.2f%s", float64(bytes)/float64(div), units[exp])
}
