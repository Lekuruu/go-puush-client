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
	checkButton := widget.NewButton("Check for Updates", func() {
		ui.RequestUpdateCheck()
	})
	updateChannel := container.NewGridWithColumns(2, branchSelect, checkButton)

	lastChecked := "Never"
	if !ui.config.Misc.LastUpdate.IsZero() {
		lastChecked = ui.config.Misc.LastUpdate.Local().Format(time.DateTime)
	}
	updateInformation := widget.NewForm(
		widget.NewFormItem("Last Checked:", widget.NewLabel(lastChecked)),
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
