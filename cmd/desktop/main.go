package main

import (
	"fyne.io/fyne/v2/app"

	"github.com/Lekuruu/go-puush-client/internal/config"
	"github.com/Lekuruu/go-puush-client/internal/desktop"
)

func main() {
	app := app.NewWithID("me.puush.client")

	store := config.NewStore()
	cfg, err := store.Load()
	if err != nil {
		// Use default config if none was found
		cfg = config.DefaultConfig()
	}

	// Save config once app shuts down
	defer store.Save(cfg)

	ui := desktop.NewUI(app, cfg)
	ui.Run()
}
