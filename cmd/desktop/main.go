package main

import (
	"fyne.io/fyne/v2/app"

	"github.com/Lekuruu/go-puush-client/internal/desktop"
)

func main() {
	app := app.NewWithID("me.puush.client")
	ui := desktop.NewUI(app)
	ui.Run()
}
