package desktop

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Lekuruu/go-puush-client/internal/updater"
)

func (ui *UI) buildUpdateTab() fyne.CanvasObject {
	autoUpdateCheckbox := widget.NewCheck("Automatically check for updates", func(enabled bool) {
		ui.config.General.AutoUpdate = enabled
	})
	autoUpdateCheckbox.Checked = ui.config.General.AutoUpdate

	branchSelect := widget.NewSelect(
		[]string{
			updater.BranchStable.String(),
			updater.BranchNightly.String(),
		},
		func(selected string) {
			ui.config.General.UpdateBranch = updater.NewBranchFromString(selected)
		},
	)
	branchSelect.SetSelected(ui.config.General.UpdateBranch.String())

	var checkButton *widget.Button
	checkButton = widget.NewButton("Check for Updates", func() {
		checkButton.SetText("Checking...")
		checkButton.Disable()
		if !ui.RequestUpdateCheck() {
			checkButton.SetText("Check for Updates")
			checkButton.Enable()
			ui.ShowNotification("Update check in progress", "Another update check is already running.")
		}
	})
	updateChannel := container.NewGridWithColumns(2, branchSelect, checkButton)

	lastCheckedLabel := widget.NewLabel(formatLastUpdateCheck(ui.config.Misc.LastUpdate))
	ui.SetUpdateFinishedCallback(func(checkedAt time.Time) {
		checkButton.SetText("Check for Updates")
		checkButton.Enable()
		lastCheckedLabel.SetText(formatLastUpdateCheck(checkedAt))
	})
	updateInformation := widget.NewForm(
		widget.NewFormItem("Last Checked:", lastCheckedLabel),
	)
	updateManagement := container.NewGridWithColumns(2, autoUpdateCheckbox, updateInformation)

	return container.NewVBox(
		widget.NewSeparator(),
		createGroup("Update Management", updateManagement),
		widget.NewSeparator(),
		createGroup("Update Channel", updateChannel),
		widget.NewSeparator(),
	)
}

func formatLastUpdateCheck(checkedAt time.Time) string {
	if checkedAt.IsZero() {
		return "Never"
	}
	return checkedAt.Local().Format(time.DateTime)
}
