/*
 * Portions of this file are derived from work by Martin Lindhe
 * licensed under the MIT License, see
 * https://github.com/martinlindhe/notify
 */

package notifications

import (
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"

	gosxnotifier "github.com/deckarep/gosx-notifier"
)

const puushBundleIdentifier = "me.puush.client"

func (n *Notification) Push() error {
	note := notification(n.Application, n.Title, n.Text, n.iconPath, n.actionUrl)
	note.Sender = puushBundleIdentifier

	if n.soundPath != "" {
		note.Sound = gosxnotifier.Sound(n.soundPath)
	}

	args, err := notificationArgs(note)
	if err != nil {
		return err
	}

	_, err = exec.Command(gosxnotifier.FinalPath, args...).Output()
	return err
}

func notificationArgs(note *gosxnotifier.Notification) ([]string, error) {
	if note.Message == "" {
		return nil, fmt.Errorf("notification message cannot be empty")
	}

	args := []string{"-message", note.Message}
	if note.Title != "" {
		args = append(args, "-title", note.Title)
	}
	if note.Subtitle != "" {
		args = append(args, "-subtitle", note.Subtitle)
	}
	if note.Sound != "" {
		args = append(args, "-sound", string(note.Sound))
	}
	if note.Group != "" {
		args = append(args, "-group", note.Group)
	}

	if note.AppIcon != "" {
		iconPath, err := filepath.Abs(note.AppIcon)
		if err != nil {
			return nil, fmt.Errorf("resolve notification icon: %w", err)
		}
		args = append(args, "-appIcon", iconPath)
	}
	if note.ContentImage != "" {
		imagePath, err := filepath.Abs(note.ContentImage)
		if err != nil {
			return nil, fmt.Errorf("resolve notification image: %w", err)
		}
		args = append(args, "-contentImage", imagePath)
	}
	if note.Link != "" {
		parsedUrl, err := url.Parse(note.Link)
		if err != nil || parsedUrl == nil {
			return nil, fmt.Errorf("invalid notification link %q", note.Link)
		}
		args = append(args, "-open", note.Link)
	}

	if strings.TrimSpace(note.Sender) != "" {
		args = append(args, "-sender", note.Sender)
	}
	return args, nil
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
