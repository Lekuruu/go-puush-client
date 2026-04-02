package screenshots

import (
	"errors"
	"io"
)

type ScreenshotProvider interface {
	// Name returns the name of the screenshot provider
	Name() string

	// Available checks if the provider is available on the system
	Available() bool

	// CaptureScreen captures the entire screen
	CaptureScreen() (io.ReadSeekCloser, error)

	// CaptureArea captures a specific region of the screen
	CaptureArea() (io.ReadSeekCloser, error)

	// CaptureWindow captures a specific window
	CaptureWindow() (io.ReadSeekCloser, error)
}

// ScreenshotProviders is a list of functions that return available screenshot providers
var ScreenshotProviders []func() (ScreenshotProvider, error)

// GetDefaultProvider returns the first available screenshot provider
func GetDefaultProvider() (ScreenshotProvider, error) {
	for _, providerFunc := range ScreenshotProviders {
		provider, err := providerFunc()
		if err == nil && provider.Available() {
			return provider, nil
		}
	}
	return nil, errors.New("no available screenshot provider found")
}
