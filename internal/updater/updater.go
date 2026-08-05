package updater

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const updateCheckTimeout = 8 * time.Second

func CanUpdate() bool {
	executable, err := os.Executable()
	if err != nil {
		return false
	}
	executablePath := filepath.Dir(executable)
	return hasWritePermission(executablePath)
}

func Check(current Version) (ReleaseCandidate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), updateCheckTimeout)
	defer cancel()

	release, err := FetchGitHubRelease(ctx)
	if err != nil {
		return nil, err
	}
	if release == nil {
		// No release found
		return nil, nil
	}
	if release.DownloadUrl() == "" {
		// No download url found
		return nil, nil
	}

	releaseVersion, err := NewVersionFromString(release.Version())
	if err != nil {
		return nil, err
	}

	if current.IsOlderThan(releaseVersion) {
		// We have found a newer version on github
		return release, nil
	}

	// No newer version found
	return nil, nil
}

func Perform(candidate ReleaseCandidate) (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("get executable path: %w", err)
	}
	executableInfo, err := os.Stat(executable)
	if err != nil {
		return "", fmt.Errorf("stat current executable: %w", err)
	}

	replacementExecutable := executable + ".new"
	defer os.Remove(replacementExecutable)

	err = downloadFile(candidate.DownloadUrl(), replacementExecutable)
	if err != nil {
		return "", fmt.Errorf("download update: %w", err)
	}

	// Preserve the mode of the current executable so the replacement remains executable on unix
	if err := os.Chmod(replacementExecutable, executableInfo.Mode().Perm()); err != nil {
		return "", fmt.Errorf("set replacement executable permissions: %w", err)
	}

	// Rename current executable to a backup name
	backupExecutable := executable + ".old"
	err = os.Rename(executable, backupExecutable)
	if err != nil {
		return "", fmt.Errorf("rename current executable: %w", err)
	}

	// Rename the new executable to the original name
	err = os.Rename(replacementExecutable, executable)
	if err != nil {
		// We're in a really bad state right now, since no executable would be available to run
		// Attempt to restore the original executable
		os.Rename(backupExecutable, executable)
		return "", fmt.Errorf("replace executable with new version: %w", err)
	}
	return executable, nil
}

func Cleanup() bool {
	executable, err := os.Executable()
	if err != nil {
		return false
	}
	newExecutable := executable + ".new"
	backupExecutable := executable + ".old"

	// Remove any leftover new/replacement executable
	// We don't really care if it fails
	os.Remove(newExecutable)

	if _, err := os.Stat(backupExecutable); err == nil {
		err = os.Remove(backupExecutable)
		if err != nil {
			return false
		}
		// An old executable was found & successfully removed
		return true
	}
	return false
}

func downloadFile(url string, destination string) error {
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	// TODO: add some context timeout here too, i'm literally too lazy right now
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func hasWritePermission(path string) bool {
	file, err := os.CreateTemp(path, ".puush-update-test-*")
	if err != nil {
		return false
	}
	originalPath := file.Name()
	renamedPath := originalPath + ".renamed"
	defer os.Remove(originalPath)
	defer os.Remove(renamedPath)

	if err := file.Close(); err != nil {
		return false
	}
	if err := os.Rename(originalPath, renamedPath); err != nil {
		return false
	}
	return os.Remove(renamedPath) == nil
}
