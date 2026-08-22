package main

import (
	"log"
	"math"
	"time"

	"github.com/Lekuruu/go-puush-client/internal/config"
	"github.com/Lekuruu/go-puush-client/internal/desktop"
	"github.com/Lekuruu/go-puush-client/internal/updater"
)

func updaterLoop(cfg *config.Config, ui *desktop.UI) {
	// We don't want to run the updater in dev builds
	if desktop.IsDevelopmentBuild() {
		return
	}

	// Check if we can even update, and inform the user if not
	if !updater.CanUpdate() {
		log.Printf("puush cannot update itself, no write permission to the executable path")
		ui.ShowNotification("puush cannot update itself!", "puush has no write permission to the installation folder. Please move it somewhere else!")
		return
	}
	currentVersion, _ := updater.NewVersionFromString(AppVersion)

	// Check if an old version is next to the current executable
	wasUpdated := updater.Cleanup()
	if wasUpdated {
		log.Printf("puush was updated to version %s", currentVersion)
		ui.ShowNotification("puush was updated!", "You are now running "+currentVersion.String()+".")
		// TODO: Add button for opening changelog page on github
	}

	// Periodically check for updates, if auto-updates are enabled inside the config
	for {
		interrupt := updaterTick(currentVersion, cfg, ui)
		if interrupt {
			return
		}

		// Start at 30 minutes and keep doubling time until we hit 6 hours max
		// This is the same thing that the official puush client did
		interval := math.Max(1000*60*30, math.Min(1000*60*360, float64(time.Hour)*2))
		intervalDuration := time.Duration(interval) * time.Millisecond
		time.Sleep(intervalDuration)
	}
}

func updaterTick(currentVersion updater.Version, cfg *config.Config, ui *desktop.UI) bool {
	if !cfg.General.AutoUpdate {
		return false
	}
	cfg.Misc.LastUpdate = time.Now()

	candidate, err := updater.Check(currentVersion, false)
	if err != nil {
		log.Printf("Failed to check for updates: %v", err)
		ui.ShowNotification("Update check failed!", "You may have to check for updates manually :(")
		return false
	}
	if candidate == nil {
		log.Printf("No update available, current version: %s", currentVersion)
		return false
	}

	log.Printf("Update available: %s -> %s", currentVersion, candidate.Version())
	ui.ShowNotification("Downloading update...", "puush will automatically restart when done!")

	updatedExecutable, err := updater.Perform(candidate)
	if err != nil {
		ui.ShowNotification("Update failed!", "You may have to install this update manually :(")
		log.Printf("Failed to perform update: %v", err)
		return false
	}

	// Stop the IPC server to let the new process start its own server
	ui.CloseIPCServer()

	log.Printf("Update downloaded, restarting...")
	if err := restart(updatedExecutable); err != nil {
		ui.ShowNotification("Update restart failed!", "Please restart puush manually to finish the update.")
		log.Printf("Failed to restart after update: %v", err)
		return false
	}
	ui.Quit()
	return true
}
