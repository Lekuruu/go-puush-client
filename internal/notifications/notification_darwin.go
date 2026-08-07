/*
 * Portions of this file are derived from work by Martin Lindhe
 * licensed under the MIT License, see
 * https://github.com/martinlindhe/notify
 */

package notifications

import (
	"errors"

	"fyne.io/fyne/v2"
)

func (n *Notification) Push() error {
	app := fyne.CurrentApp()
	if app == nil {
		return errors.New("app is not initialized")
	}

	content := n.Text
	if n.Title != "" {
		if content == "" {
			content = n.Title
		} else {
			content = n.Title + "\n" + content
		}
	}
	app.SendNotification(fyne.NewNotification(n.Application, content))
	return nil
}
