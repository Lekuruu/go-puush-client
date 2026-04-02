//go:build windows

package screenshots

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"time"
)

type WindowsScreenshotProvider struct{}

func NewWindowsScreenshotProvider() (ScreenshotProvider, error) {
	p := &WindowsScreenshotProvider{}
	if !p.Available() {
		return nil, errors.New("windows screenshot provider is not available")
	}
	return p, nil
}

func (p *WindowsScreenshotProvider) Name() string {
	return "Windows (Native)"
}

func (p *WindowsScreenshotProvider) Available() bool {
	return true // TODO
}

func (p *WindowsScreenshotProvider) CaptureScreen() (io.ReadSeekCloser, error) {
	x := getSystemMetrics(smXVirtualScreen)
	y := getSystemMetrics(smYVirtualScreen)
	width := getSystemMetrics(smCXVirtualScreen)
	height := getSystemMetrics(smCYVirtualScreen)

	img, err := captureScreenRect(x, y, width, height)
	if err != nil {
		return nil, fmt.Errorf("capture screen: %w", err)
	}

	return newPngReader(img)
}

func (p *WindowsScreenshotProvider) CaptureArea() (io.ReadSeekCloser, error) {
	r, err := selectAreaRect()
	if err != nil {
		return nil, fmt.Errorf("select area: %w", err)
	}

	top := int(r.Top)
	left := int(r.Left)
	width := int(r.Right - r.Left)
	height := int(r.Bottom - r.Top)

	// Give the selector window a moment to disappear before capturing
	time.Sleep(50 * time.Millisecond)

	img, err := captureScreenRect(left, top, width, height)
	if err != nil {
		return nil, fmt.Errorf("capture area: %w", err)
	}

	return newPngReader(img)
}

func (p *WindowsScreenshotProvider) CaptureWindow() (io.ReadSeekCloser, error) {
	hwnd, err := getForegroundWindow()
	if err != nil {
		return nil, fmt.Errorf("get foreground window: %w", err)
	}

	img, err := captureWindow(hwnd)
	if err != nil {
		return nil, fmt.Errorf("capture window: %w", err)
	}

	return newPngReader(img)
}

func newPngReader(img image.Image) (io.ReadSeekCloser, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("encode png: %w", err)
	}
	return &memoryReadCloser{bytes.NewReader(buf.Bytes())}, nil
}

func init() {
	ScreenshotProviders = append(ScreenshotProviders, NewWindowsScreenshotProvider)
}
