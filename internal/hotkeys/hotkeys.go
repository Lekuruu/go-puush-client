package hotkeys

import (
	"log"
	"strings"

	"github.com/Lekuruu/go-hotkey"
	"github.com/Lekuruu/go-puush-client/internal/config"
	"github.com/Lekuruu/go-puush-client/internal/tray"
)

// HotkeyManager handles global hotkeys safely across Start/Stop cycles.
type HotkeyManager struct {
	config  *config.Config
	tray    *tray.TrayManager
	hotkeys map[string]*hotkey.Hotkey
}

func NewHotkeyManager(cfg *config.Config, tm *tray.TrayManager) *HotkeyManager {
	return &HotkeyManager{
		tray:    tm,
		config:  cfg,
		hotkeys: make(map[string]*hotkey.Hotkey),
	}
}

// Start registers all hotkeys defined in the configuration.
func (m *HotkeyManager) Start() {
	m.register(m.config.Hotkeys.ScreenSelection, m.tray.UploadAreaScreenshot)
	m.register(m.config.Hotkeys.FullscreenScreenshot, m.tray.UploadDesktopScreenshot)
	m.register(m.config.Hotkeys.CurrentWindowScreenshot, m.tray.UploadWindowScreenshot)
	m.register(m.config.Hotkeys.UploadFile, m.tray.UploadFileFromDialog)
	m.register(m.config.Hotkeys.UploadClipboard, m.tray.UploadFromClipboard)
	m.register(m.config.Hotkeys.Toggle, m.tray.TogglePuushing)
}

// Stop unregisters all active hotkeys and clears the hotkey map.
func (m *HotkeyManager) Stop() {
	for shortcut, hk := range m.hotkeys {
		if err := hk.Unregister(); err != nil {
			log.Printf("Failed to unregister hotkey %s: %v", shortcut, err)
		}
	}
	m.hotkeys = make(map[string]*hotkey.Hotkey)
	log.Println("All hotkeys unregistered")
}

func (m *HotkeyManager) register(shortcut string, action func()) {
	if shortcut == "" {
		return
	}

	parts := strings.Split(shortcut, "+")
	if len(parts) < 2 {
		log.Printf("Invalid hotkey format: %s", shortcut)
		return
	}

	mods := parseModifiers(parts)
	key := parseKey(parts[len(parts)-1])

	hk := hotkey.New(mods, key)
	err := hk.Register()

	if err != nil {
		log.Printf("Failed to register hotkey %s: %v", shortcut, err)
		return
	}
	m.hotkeys[shortcut] = hk

	// Start a goroutine to listen for this hotkey's events
	go func(hk *hotkey.Hotkey) {
		for {
			_, ok := <-hk.Keydown()
			if !ok {
				return
			}
			if m.tray.PuushingDisabled() {
				continue
			}
			action()
		}
	}(hk)

	log.Printf("Registered hotkey: %s (%v)", shortcut, hk)
}
