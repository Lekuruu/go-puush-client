package screenshots

import "io"

type ScreenshotProvider interface {
	// Name returns the name of the screenshot provider
	Name() string

	// Available checks if the provider is available on the system
	Available() bool

	// CaptureScreen captures the entire screen
	CaptureScreen() (io.ReadCloser, error)

	// CaptureArea captures a specific region of the screen
	CaptureArea() (io.ReadCloser, error)

	// CaptureWindow captures a specific window
	CaptureWindow() (io.ReadCloser, error)
}

// ScreenshotProviders is a list of functions that return available screenshot providers
var ScreenshotProviders []func() (ScreenshotProvider, error)
