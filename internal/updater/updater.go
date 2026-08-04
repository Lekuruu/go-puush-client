package updater

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func Check(current Version) (ReleaseCandidate, error) {
	release, err := FetchGitHubRelease(context.Background()) // TODO: add context with timeout
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

func Perform(candidate ReleaseCandidate) error {
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get executable path: %w", err)
	}

	executablePath := filepath.Dir(executable)
	if !hasWritePermission(executablePath) {
		return fmt.Errorf("no write permission to the executable path: %s", executablePath)
	}

	replacementExecutable := executable + ".new"
	err = downloadFile(candidate.DownloadUrl(), replacementExecutable)
	if err != nil {
		return fmt.Errorf("download update: %w", err)
	}

	// Rename current executable to a backup name
	backupExecutable := executable + ".old"
	err = os.Rename(executable, backupExecutable)
	if err != nil {
		return fmt.Errorf("rename current executable: %w", err)
	}

	// Rename the new executable to the original name
	err = os.Rename(replacementExecutable, executable)
	if err != nil {
		// We're in a really bad state right now, since no executable would be available to run
		// Attempt to restore the original executable
		os.Rename(backupExecutable, executable)
		return fmt.Errorf("replace executable with new version: %w", err)
	}
	return nil
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

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func hasWritePermission(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	mode := fileInfo.Mode()
	if mode&os.ModePerm == os.ModePerm {
		return false
	} else {
		return true
	}
}
