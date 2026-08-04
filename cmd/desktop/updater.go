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
		// Start at 30 minutes and keep doubling time until we hit 6 hours max
		// This is the same thing that the official puush client did
		interval := math.Max(1000*60*30, math.Min(1000*60*360, float64(time.Hour)*2))
		intervalDuration := time.Duration(interval) * time.Millisecond

		if !cfg.Misc.UpdatesEnabled {
			time.Sleep(intervalDuration)
			return
		}

		candidate, err := updater.Check(currentVersion)
		if err != nil {
			log.Printf("Failed to check for updates: %v", err)
			ui.ShowNotification("Update check failed!", "You may have to check for updates manually :(")
			time.Sleep(intervalDuration)
			continue
		}
		if candidate == nil {
			log.Printf("No update available, current version: %s", currentVersion)
			time.Sleep(intervalDuration)
			continue
		}

		log.Printf("Update available: %s -> %s", currentVersion, candidate.Version())
		ui.ShowNotification("Downloading update...", "puush will automatically restart when done!")

		err = updater.Perform(candidate)
		if err != nil {
			ui.ShowNotification("Update failed!", "You may have to install this update manually :(")
			log.Printf("Failed to perform update: %v", err)
			time.Sleep(intervalDuration)
			continue
		}

		log.Printf("Update downloaded, restarting...")
		restart()
	}
}
