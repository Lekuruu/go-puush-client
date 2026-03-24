/*
 * Portions of this file are derived from work by Martin Lindhe
 * licensed under the MIT License, see
 * https://github.com/martinlindhe/notify
 */

package notifications

import (
	gosxnotifier "github.com/deckarep/gosx-notifier"
)

func (n *Notification) Push() error {
	note := notification(n.Application, n.Title, n.Text, n.iconPath, n.actionUrl)
	if n.soundPath != "" {
		note.Sound = gosxnotifier.Sound(n.soundPath)
	}
	return note.Push()
}

func notification(appName string, title string, text string, iconPath string, actionUrl string) *gosxnotifier.Notification {
	head := ""
	if text == "" {
		head = title
		title = ""
	} else {
		head = text
	}

	note := gosxnotifier.NewNotification(head)
	note.Title = appName
	note.Subtitle = title
	note.AppIcon = iconPath // (10.9+ ONLY)
	if actionUrl != "" {
		note.Link = actionUrl
	}
	return note
}
