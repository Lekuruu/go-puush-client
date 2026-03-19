package main

import (
	"fyne.io/fyne/v2/app"

	"github.com/Lekuruu/go-puush-client/internal/config"
	"github.com/Lekuruu/go-puush-client/internal/desktop"
	"github.com/Lekuruu/go-puush-client/pkg/puush"
)

func main() {
	app := app.NewWithID("me.puush.client")

	store := config.NewStore()
	cfg, err := store.Load()
	if err != nil {
		// Use default config if none was found
		cfg = config.DefaultConfig()
	}

	api := puush.NewClientFromApiKey(cfg.Account.Username, cfg.Account.Key)
	defer func() {
		if !api.Account.Credentials.HasApiKey() {
			return
		}
		// Update credentials in config after shutdown
		cfg.Account.Key = *api.Account.Credentials.Key
		cfg.Account.Username = *api.Account.Credentials.Identifier
	}()

	// Save config once app shuts down
	defer store.Save(cfg)

	ui := desktop.NewUI(app, api, cfg)
	ui.Run()
}
