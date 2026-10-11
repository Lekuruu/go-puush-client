//go:build linux

package screenshots

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

// https://github.com/Netflate/LumineCapture

type LumineScreenshotProvider struct {
	binPath        string
	timeout        time.Duration
	fullscreenMode FullscreenMode
	quality        Quality
}

func NewLumineCaptureProvider() (ScreenshotProvider, error) {
	binPath, err := exec.LookPath("lumine-capture")
	if err != nil {
		return nil, fmt.Errorf("lumine-capture not found in PATH: %w", err)
	}

	return &LumineScreenshotProvider{
		binPath: binPath,
		timeout: 5 * time.Minute,
	}, nil
}

func (p *LumineScreenshotProvider) Name() string {
	return "LumineCapture"
}

func (p *LumineScreenshotProvider) Warning() string {
	return "LumineCapture is wayland-only, and is recommended to be used with KDE Plasma."
}

func (p *LumineScreenshotProvider) SetQuality(quality Quality) {
	p.quality = quality
}

func (p *LumineScreenshotProvider) SetFullscreenMode(mode FullscreenMode) {
	p.fullscreenMode = mode
}

func (p *LumineScreenshotProvider) Available() bool {
	_, err := exec.LookPath("lumine-capture")
	return err == nil
}

// CaptureScreen captures either the whole workspace or the monitor under the pointer
func (p *LumineScreenshotProvider) CaptureScreen() (io.ReadSeekCloser, error) {
	args := []string{"-f"}
	if p.fullscreenMode != FullscreenModeAllScreens {
		args = append(args, "-m")
	}
	return p.performCapture(args...)
}

// CaptureArea lets the user drag a region to capture
func (p *LumineScreenshotProvider) CaptureArea() (io.ReadSeekCloser, error) {
	return p.performCapture("-r")
}

// CaptureWindow captures the active window (KDE Plasma only)
func (p *LumineScreenshotProvider) CaptureWindow() (io.ReadSeekCloser, error) {
	return p.performCapture("-w")
}

func (p *LumineScreenshotProvider) performCapture(modeArgs ...string) (io.ReadSeekCloser, error) {
	dir, err := os.MkdirTemp("", "lumine-capture-*")
	if err != nil {
		return nil, fmt.Errorf("create temp directory: %w", err)
	}

	cleanup := true
	defer func() {
		if cleanup {
			os.RemoveAll(dir)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	args := append(modeArgs, "-t", "s", "-o", dir)
	cmd := exec.CommandContext(ctx, p.binPath, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("lumine-capture timed out")
		}

		msg := strings.TrimSpace(string(output))
		if msg != "" {
			return nil, fmt.Errorf("lumine-capture failed: %w: %s", err, msg)
		}
		return nil, fmt.Errorf("lumine-capture failed or was cancelled: %w", err)
	}

	path, err := findScreenshotInDirectory(dir)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open screenshot: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("stat screenshot: %w", err)
	}
	if info.Size() == 0 {
		file.Close()
		return nil, errors.New("screenshot was cancelled or is empty")
	}

	cleanup = false
	reader := &temporaryDirectoryReader{ // will remove the directory on Close()
		file: file,
		dir:  dir,
	}
	return ApplyQuality(reader, p.quality)
}

func init() {
	ScreenshotProviders = append(ScreenshotProviders, NewLumineCaptureProvider)
}
