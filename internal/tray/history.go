package tray

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"github.com/Lekuruu/go-puush-client/pkg/puush"
)

// RefreshHistory will update the tray's upload history
func (m *TrayManager) RefreshHistory() {
	if !m.api.Account.Credentials.HasApiKey() {
		return
	}
	history, err := m.api.History()
	if err != nil {
		return
	}
	m.uploadHistory = history
	m.rebuildMenuItems()
}

func (m *TrayManager) BuildHistoryMenu() []*fyne.MenuItem {
	recentUploads := fyne.NewMenuItem("Recent Uploads", func() {})
	recentUploads.Disabled = true
	items := []*fyne.MenuItem{recentUploads}

	for _, historyItem := range m.uploadHistory {
		items = append(items, m.BuildHistoryMenuItem(historyItem))
	}
	return items
}

func (m *TrayManager) BuildHistoryMenuItem(historyItem *puush.HistoryItem) *fyne.MenuItem {
	timeItem := fyne.NewMenuItem(fmt.Sprintf("Uploaded: %s", historyItem.Time.Format("2006-01-02 15:04:05")), func() {})
	timeItem.Disabled = true

	viewsItem := fyne.NewMenuItem(fmt.Sprintf("Views: %d", historyItem.Views), func() {})
	viewsItem.Disabled = true

	openItem := fyne.NewMenuItem("Open in browser", func() {
		if u, err := url.Parse(historyItem.Url); err == nil {
			fyne.CurrentApp().OpenURL(u)
		}
	})

	copyItem := fyne.NewMenuItem("Copy link to clipboard", func() {
		fyne.CurrentApp().Clipboard().SetContent(historyItem.Url)
	})

	deleteItem := fyne.NewMenuItem("Delete", func() {
		if newHistory, err := m.api.Delete(historyItem.Id); err == nil {
			m.uploadHistory = newHistory
			m.rebuildMenuItems()
		}
	})

	historyMenu := fyne.NewMenu(historyItem.FileName,
		timeItem,
		viewsItem,
		fyne.NewMenuItemSeparator(),
		openItem,
		copyItem,
		fyne.NewMenuItemSeparator(),
		deleteItem,
	)

	// TODO: Add icons to history menu items
	historyMenuItem := fyne.NewMenuItem(historyItem.FileName, nil)
	historyMenuItem.ChildMenu = historyMenu
	return historyMenuItem
}
