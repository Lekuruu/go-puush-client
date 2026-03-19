package desktop

import (
	"image/color"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/Lekuruu/go-puush-client/assets"
)

// puush uses a pre-made background asset for the quick start window.
// Elements like the login prompt are just added on top of it.

func (ui *UI) ShowStartupWindow() {
	w := ui.app.NewWindow("puush quick start")
	w.SetFixedSize(true)

	// Create the background image from our embedded asset
	bgResource := fyne.NewStaticResource("quickstart_bg", assets.QuickstartData)
	bgImage := canvas.NewImageFromResource(bgResource)
	bgImage.FillMode = canvas.ImageFillCover
	bgImage.SetMinSize(fyne.NewSize(640, 540))

	// Create button to link to account page
	registerBtn := NewBorderedButton("Take me to the account creation page!", func() {
		OpenBrowser("https://puush.me/register") // TODO: Custom server url
	})
	registerBtn.Move(fyne.NewPos(200, 138))
	registerBtn.Resize(fyne.NewSize(250, 28))

	emailLabel := canvas.NewText("Email:", color.Black)
	emailLabel.Move(fyne.NewPos(155, 207))

	passwordLabel := canvas.NewText("Password:", color.Black)
	passwordLabel.Move(fyne.NewPos(132, 237))

	// Create the inputs that will be placed over the background
	emailEntry := widget.NewEntry()
	emailEntry.Move(fyne.NewPos(200, 200))
	emailEntry.Resize(fyne.NewSize(150, 25))

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.Move(fyne.NewPos(200, 230))
	passwordEntry.Resize(fyne.NewSize(150, 25))

	// Forgot password hyperlink
	forgotURL, _ := url.Parse("https://puush.me/reset_password") // TODO: Custom server url
	forgotLink := NewUnderlinedLink("Forgotten Password?", forgotURL)
	forgotLink.Move(fyne.NewPos(201, 256))
	forgotLink.Resize(forgotLink.MinSize())

	// Create a rectangle to cover the "Login successful" text
	coverRectangle := canvas.NewRectangle(color.White)
	coverRectangle.Move(fyne.NewPos(125, 200))
	coverRectangle.Resize(fyne.NewSize(400, 75))

	loginBtn := NewBorderedButton("Login", func() {
		// TODO: Implement login logic
		// When a login is successful, the username, password & cover rect should be hidden
		// When a login failed, an error message box should pop up with the appropriate error message
	})
	loginBtn.Move(fyne.NewPos(370, 200))
	loginBtn.Resize(fyne.NewSize(160, 55))
	loginBtn.Instance.Disable()

	// Container for the background and absolutely positioned overlays
	overlayContainer := container.NewWithoutLayout(
		bgImage,
		coverRectangle,
		registerBtn,
		emailLabel,
		passwordLabel,
		emailEntry,
		passwordEntry,
		forgotLink,
		loginBtn,
	)
	bgContainer := container.NewStack(bgImage, overlayContainer)

	startupCheckbox := widget.NewCheck("Start puush on startup", func(checked bool) {
		// TODO: Implement startup logic
	})
	startupCheckbox.SetChecked(true)

	okayBtn := NewBorderedButton("Okay, I've got it!", func() {
		w.Close()
	})
	okayBtn.Instance.Disable()

	// Force-resize this button inside a new container, since the HBox will not allow that
	sizedOkayBtn := container.NewGridWrap(fyne.NewSize(325, 32), okayBtn)

	// Layout the bottom bar
	bottomBar := container.NewHBox(
		layout.NewSpacer(),
		startupCheckbox,
		layout.NewSpacer(),
		sizedOkayBtn,
		layout.NewSpacer(),
	)

	// Combine the background area and the bottom bar
	mainContent := container.NewVBox(
		bgContainer,
		widget.NewSeparator(),
		container.NewPadded(bottomBar),
	)
	w.SetContent(mainContent)
	w.Resize(fyne.NewSize(640, 540))

	// `SetMaster` will close the process once the window is closed
	// TODO: Disable `SetMaster` once we have the tray system running
	w.SetMaster()
	w.Show()
}
