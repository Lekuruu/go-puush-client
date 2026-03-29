//go:build linux

package screenshots

import (
	"testing"
)

func TestSpectacleProvider(t *testing.T) {
	provider, err := NewSpectacleProvider()
	if err != nil {
		t.Skipf("failed to create spectacle provider: %v", err)
	}

	t.Run("CaptureScreen", func(t *testing.T) {
		reader, err := provider.CaptureScreen()
		if err != nil {
			t.Fatalf("CaptureScreen failed: %v", err)
		}
		defer reader.Close()
	})

	t.Run("CaptureArea", func(t *testing.T) {
		reader, err := provider.CaptureArea()
		if err != nil {
			t.Fatalf("CaptureArea failed: %v", err)
		}
		defer reader.Close()
	})

	t.Run("CaptureWindow", func(t *testing.T) {
		reader, err := provider.CaptureWindow()
		if err != nil {
			t.Fatalf("CaptureWindow failed: %v", err)
		}
		defer reader.Close()
	})
}
