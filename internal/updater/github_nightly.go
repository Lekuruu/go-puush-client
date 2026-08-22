package updater

import (
	"context"
	"runtime"
	"time"

	"github.com/google/go-github/v89/github"
)

const GitHubWorkflowName = "build.yml"

type GitHubNightlyCandidate struct {
	workflowRun       *github.WorkflowRun
	workflowArtifacts []*github.Artifact
}

func (c *GitHubNightlyCandidate) Version() Version {
	commit := c.workflowRun.GetRunNumber()
	return NewIntegerVersion(commit)
}

func (c *GitHubNightlyCandidate) Branch() Branch {
	return BranchNightly
}

func (c *GitHubNightlyCandidate) Description() string {
	return c.workflowRun.GetHeadCommit().GetMessage()
}

func (c *GitHubNightlyCandidate) CreatedAt() time.Time {
	return c.workflowRun.GetCreatedAt().Time
}

func (c *GitHubNightlyCandidate) IsPrerelease() bool {
	return true
}

func (c *GitHubNightlyCandidate) DownloadUrl() string {
	targetFilename := releaseAssetFilename(
		runtime.GOOS,
		runtime.GOARCH,
	)
	for _, artifact := range c.workflowArtifacts {
		if artifact.GetName() == targetFilename {
			return artifact.GetArchiveDownloadURL()
		}
	}
	return ""
}

func FetchGitHubNightly(ctx context.Context) (ReleaseCandidate, error) {
	client, err := github.NewClient()
	if err != nil {
		return nil, err
	}

	workflowRuns, _, err := client.Actions.ListWorkflowRunsByFileName(ctx, GitHubUser, GitHubRepository, GitHubWorkflowName, &github.ListWorkflowRunsOptions{
		Status: "completed",
	})
	if err != nil {
		return nil, err
	}
	if len(workflowRuns.WorkflowRuns) == 0 {
		return nil, nil
	}

	latestRun := workflowRuns.WorkflowRuns[0]
	artifacts, _, err := client.Actions.ListWorkflowRunArtifacts(ctx, GitHubUser, GitHubRepository, latestRun.GetID(), &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	return &GitHubNightlyCandidate{
		workflowRun:       latestRun,
		workflowArtifacts: artifacts.Artifacts,
	}, nil
}
