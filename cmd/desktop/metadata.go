package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Lekuruu/go-puush-client/assets"
)

const (
	AppName     = "puush"
	AppID       = "me.puush.client"
	AppVersion  = "1.0.0"
	AppBuild    = 100
	AppIconName = "puush.png"
)

func init() {
	app.SetMetadata(fyne.AppMetadata{
		ID:      AppID,
		Name:    AppName,
		Version: AppVersion,
		Build:   AppBuild,
		Icon: &fyne.StaticResource{
			StaticName:    AppIconName,
			StaticContent: assets.PuushIconData,
		},
		Release:    true,
		Custom:     map[string]string{},
		Migrations: map[string]bool{"fyneDo": true},
	})
}
