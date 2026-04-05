package desktop

import (
	"fmt"
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
	emailLabel := trailingLabel("Email:")
	passwordLabel := trailingLabel("Password:")

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
		// TODO: Disable buttons when performing login & do it inside a go routine
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
	accountTypeStr := ui.config.Account.Type.String() + " Account"
	diskUsageStr := ui.config.Account.DiskUsageHumanReadable()

	expiryStr := ui.config.Account.Expiry
	if expiryStr == "" {
		expiryStr = "Never"
	}

	detailsGrid := container.NewGridWithColumns(2,
		trailingLabel("Logged in as:"), widget.NewLabel(ui.config.Account.Username),
		trailingLabel("API Key:"), widget.NewLabel(ui.config.Account.Key),
		trailingLabel("Account Type:"), widget.NewLabel(accountTypeStr),
		trailingLabel("Expiry Date:"), widget.NewLabel(expiryStr),
		trailingLabel("Disk Usage:"), widget.NewLabel(diskUsageStr),
	)

	myAccountBtn := widget.NewButton("My Account", func() {
		path := fmt.Sprintf("/login/go/?k=%s", ui.config.Account.Key)
		OpenBrowser(ui.api.FormatURL(path))
	})

	logoutBtn := widget.NewButton("Logout", func() {
		ui.config.Account.Reset()
		ui.api.Account.Reset()
		updateView()
	})

	buttons := container.NewGridWithColumns(
		2, myAccountBtn, logoutBtn,
	)
	content := container.NewVBox(
		container.NewPadded(detailsGrid),
		widget.NewLabel(""),
		container.NewPadded(buttons),
	)
	return createGroup("Account Details", content)
}
