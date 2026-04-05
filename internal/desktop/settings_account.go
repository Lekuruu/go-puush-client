package desktop

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Lekuruu/go-puush-client/pkg/puush"
)

func (ui *UI) buildAccountTab() fyne.CanvasObject {
	var updateView func()

	accountContainer := container.NewStack()
	updateView = func() {
		defer accountContainer.Refresh()
		accountContainer.Objects = nil
		if ui.config.Account.Key != "" {
			accountContainer.Add(ui.buildAccountDetails(updateView))
		} else {
			accountContainer.Add(ui.buildAccountSetup(updateView))
		}
	}
	updateView()

	return container.NewVBox(
		widget.NewSeparator(),
		accountContainer,
		widget.NewSeparator(),
	)
}

func (ui *UI) buildAccountSetup(updateView func()) fyne.CanvasObject {
	infoText := "You need to login before you can make full use of puush. "
	infoText += "If you don't already have an account, you can register for free via the link below."
	infoLabel := widget.NewLabel(infoText)
	infoLabel.Wrapping = fyne.TextWrapWord

	serverUrl := ui.config.Misc.ParseServerURL()
	resetUrl := serverUrl.String() + "/reset_password"
	registerUrl := serverUrl.String() + "/register"

	emailEntry := widget.NewEntry()
	passwordEntry := widget.NewPasswordEntry()

	emailLabel := widget.NewLabel("Email:")
	emailLabel.Alignment = fyne.TextAlignTrailing
	passwordLabel := widget.NewLabel("Password:")
	passwordLabel.Alignment = fyne.TextAlignTrailing

	form := container.NewGridWithColumns(2,
		emailLabel, emailEntry,
		passwordLabel, passwordEntry,
	)

	forgotURL, _ := url.Parse(resetUrl)
	registerURL, _ := url.Parse(registerUrl)
	forgotLink := NewUnderlinedLink("Forgotten Password?", forgotURL)
	registerLink := NewUnderlinedLink("Sign up for free account...", registerURL)
	linksContainer := container.NewHBox(forgotLink, layout.NewSpacer(), registerLink)

	performLogin := func() {
		ui.api.Account.Credentials = &puush.Credentials{
			Identifier: &emailEntry.Text,
			Password:   &passwordEntry.Text,
		}
		ui.api.SetBaseURL(serverUrl.String())
		if err := ui.api.Authenticate(); err != nil {
			showError(err)
			return
		}
		ui.UpdateAccountConfiguration()
		updateView()
		// TODO: Disable buttons when performing login
	}
	loginBtn := NewBorderedButton("Login", performLogin)

	sizedForm := container.NewGridWrap(fyne.NewSize(350, 55), form)
	sizedLoginBtn := container.NewGridWrap(fyne.NewSize(140, 53), loginBtn)

	loginContainer := container.NewHBox(
		layout.NewSpacer(),
		sizedForm, widget.NewLabel(" "), sizedLoginBtn,
		layout.NewSpacer(),
	)

	content := container.NewVBox(
		infoLabel,
		widget.NewLabel(""),
		loginContainer,
		widget.NewLabel(""),
		linksContainer,
	)

	return createGroup("Account Setup", content)
}

func (ui *UI) buildAccountDetails(updateView func()) fyne.CanvasObject {
	// TODO
	return nil
}
