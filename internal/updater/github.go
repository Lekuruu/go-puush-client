package updater

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/google/go-github/v89/github"
)

const GitHubUser = "Lekuruu"
const GitHubRepository = "go-puush-client"

type GitHubReleaseCandidate struct {
	release *github.RepositoryRelease
}

func (c *GitHubReleaseCandidate) Version() string {
	return c.release.GetTagName()
}

func (c *GitHubReleaseCandidate) Branch() Branch {
	return BranchStable
}

func (c *GitHubReleaseCandidate) Description() string {
	return c.release.GetBody()
}

func (c *GitHubReleaseCandidate) CreatedAt() time.Time {
	return c.release.GetCreatedAt().Time
}

func (c *GitHubReleaseCandidate) IsPrerelease() bool {
	return c.release.GetPrerelease()
}

func (c *GitHubReleaseCandidate) DownloadUrl() string {
	targetFilename := releaseAssetFilename(
		runtime.GOOS,
		runtime.GOARCH,
	)
	for _, asset := range c.release.Assets {
		if asset.GetName() == targetFilename {
			return asset.GetBrowserDownloadURL()
		}
	}
	return ""
}

func FetchGitHubRelease(ctx context.Context) (ReleaseCandidate, error) {
	client, err := github.NewClient()
	if err != nil {
		return nil, err
	}

	// NOTE: This does not include pre-releases, so we'll always consider this as stable updates
	release, response, err := client.Repositories.GetLatestRelease(ctx, GitHubUser, GitHubRepository)
	if response.StatusCode == 404 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &GitHubReleaseCandidate{release: release}, nil
}

func releaseAssetFilename(goos, goarch string) string {
	switch goos {
	case "darwin":
		return fmt.Sprintf("puush-macos-%s.app.zip", goarch)
	case "windows":
		return fmt.Sprintf("puush-windows-%s.exe", goarch)
	default:
		return fmt.Sprintf("puush-%s-%s", goos, goarch)
	}
}
