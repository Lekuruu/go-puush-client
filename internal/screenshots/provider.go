package screenshots

import "io"

type ScreenshotProvider interface {
	// CaptureScreen captures the entire screen
	CaptureScreen() (io.ReadCloser, error)

	// CaptureArea captures a specific region of the screen
	CaptureArea() (io.ReadCloser, error)

	// CaptureWindow captures a specific window
	CaptureWindow() (io.ReadCloser, error)
}
