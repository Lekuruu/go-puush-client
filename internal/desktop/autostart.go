package desktop

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Lekuruu/go-puush-client/assets"
	"github.com/emersion/go-autostart"
)

func EnableAutostart() error {
	app, err := getAutostartApp()
	if err != nil {
		return err
	}

	if app.IsEnabled() {
		return nil
	}
	return app.Enable()
}

func DisableAutostart() error {
	app, err := getAutostartApp()
	if err != nil {
		return err
	}

	if !app.IsEnabled() {
		return nil
	}
	return app.Disable()
}

func getAutostartApp() (*autostart.App, error) {
	executable, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	// Check if this is a development build
	if filepath.Base(executable) == "desktop" && strings.Contains(executable, "go-build") {
		return nil, fmt.Errorf("autostart is not supported in development builds")
	}

	return &autostart.App{
		Name:        "puush",
		DisplayName: "puush",
		Icon:        exportIcon(),
		Exec:        []string{executable},
	}, nil
}

func exportIcon() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	iconPath := filepath.Join(cacheDir, "puush", "puush.png")

	// Check if icon already exists
	_, err = os.Stat(iconPath)
	if err == nil {
		return iconPath
	}

	// Export icon if it doesn't exist
	os.MkdirAll(filepath.Dir(iconPath), 0755)
	os.WriteFile(iconPath, assets.PuushIconData, 0644)
	return iconPath
}
