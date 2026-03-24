package main

import (
	"time"

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

	// Save config once app shuts down
	defer store.Save(cfg)

	api := puush.NewClientFromApiKey(cfg.Account.Username, cfg.Account.Key)
	api.SetBaseURL(cfg.Misc.ParseServerURL().String())
	defer func() {
		if !api.Account.Credentials.HasApiKey() {
			return
		}
		// Update account state in config after shutdown
		cfg.Account.Key = *api.Account.Credentials.Key
		cfg.Account.Username = *api.Account.Credentials.Identifier
		cfg.Account.Type = int(api.Account.Type)
		cfg.Account.Usage = api.Account.DiskUsage

		if api.Account.SubscriptionEnd != nil {
			cfg.Account.Expiry = api.Account.SubscriptionEnd.Format(time.DateTime)
		}
	}()

	// Apply previous account state from config to api
	expiry, _ := time.Parse(time.DateTime, cfg.Account.Expiry)
	api.Account.Type = puush.AccountType(cfg.Account.Type)
	api.Account.DiskUsage = cfg.Account.Usage
	api.Account.SubscriptionEnd = &expiry
	// TODO: Handle time parsing better, I'm just too lazy right now and also cba tbh

	ui := desktop.NewUI(app, api, cfg)
	ui.Run()
}
