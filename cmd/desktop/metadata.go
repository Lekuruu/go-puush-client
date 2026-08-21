package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Lekuruu/go-puush-client/assets"
)

const (
	AppName     = "puush"
	AppID       = "me.puush.client"
	AppVersion  = "1.0.4"
	AppIconName = "puush.png"
)

// Set by the build system through -ldflags -X
var (
	AppBuild  = "0"   // Number of current commits
	AppCommit = "dev" // Current commit sha
)

func init() {
	buildNumber, _ := strconv.Atoi(AppBuild)

	app.SetMetadata(fyne.AppMetadata{
		ID:      AppID,
		Name:    AppName,
		Version: AppVersion,
		Build:   buildNumber,
		Icon: &fyne.StaticResource{
			StaticName:    AppIconName,
			StaticContent: assets.PuushIconData,
		},
		Release:    true,
		Custom:     map[string]string{"commit": AppCommit},
		Migrations: map[string]bool{"fyneDo": true},
	})
}
