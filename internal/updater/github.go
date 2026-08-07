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
	extension := ""
	if goos == "windows" {
		extension = ".exe"
	}
	if goos == "darwin" {
		// Github actions publishes as "macos" instead of "darwin"
		goos = "macos"
	}
	return fmt.Sprintf("puush-%s-%s%s", goos, goarch, extension)
}
